package cmd

import (
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/spf13/cobra"
	"github.com/zebbra/dping_exporter/lib/canary"
	"github.com/zebbra/dping_exporter/lib/ping"
)

func init() {
	agentCmd.Flags().StringArrayVarP(&agentCanaryTargets, "canary.target", "c", []string{}, "addresses to check to determine if agent can reach network targets")
	agentCmd.Flags().IntVarP(&agentCanaryProbeCount, "canary.probe.count", "", agentCanaryProbeCount, "number of pings sent for each target")
	agentCmd.Flags().DurationVarP(&agentCanaryProbeInterval, "canary.probe.interval", "", agentCanaryProbeInterval, "interval for individual pings during a target check")
	agentCmd.Flags().DurationVarP(&agentCanaryProbeTimeout, "canary.probe.timeout", "", agentCanaryProbeTimeout, "timeout for a target check")
	agentCmd.Flags().IntVarP(&agentCanaryLivenessThreshold, "canary.liveness-threshold", "t", agentCanaryLivenessThreshold, "number of targets required to be alive for canary check to succeed")
	agentCmd.Flags().DurationVarP(&agentCanaryCheckInterval, "canary.interval", "i", agentCanaryCheckInterval, "how often canaries are checked for liveness")

	rootCmd.AddCommand(agentCmd)
}

var agentCanaryTargets = []string{}
var agentCanaryProbeCount = 5
var agentCanaryProbeInterval = 500 * time.Millisecond
var agentCanaryProbeTimeout = 5 * time.Second
var agentCanaryLivenessThreshold = 1
var agentCanaryCheckInterval = 5 * time.Second

var agentCmd = &cobra.Command{
	Use:   "agent [flags] <listen-address>",
	Short: "Runs polling agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		online := atomic.Bool{}
		online.Store(true)

		// activate periodic canary checks if requested
		if len(agentCanaryTargets) > 0 {
			logger := slog.With(slog.Any("targets", agentCanaryTargets))
			logger.Info(
				"starting canary check routine",
				slog.String("interval", agentCanaryCheckInterval.String()),
			)

			cp := canary.New(agentCanaryTargets)

			check := func() {
				logger.Debug("run canary probe")
				alive, _ := cp.Alive(cmd.Context())
				online.Store(alive)

				l := logger.With(slog.Bool("online", alive))
				l.Debug("probe result")

				if !alive {
					l.Error("canary check failed, offline")
				}
			}

			check()

			ticker := time.NewTicker(agentCanaryProbeTimeout)

			go func() {
				for {
					select {
					case <-cmd.Context().Done():
						return

					case <-ticker.C:
						check()
					}
				}
			}()
		}

		app := fiber.New()

		app.Get("/", func(c fiber.Ctx) error {
			return c.SendString("ready for service...")
		})

		app.Get("/-/healthy", healthcheck.New(healthcheck.Config{
			Probe: func(c fiber.Ctx) bool {
				return online.Load()
			},
		}))

		app.Get("/probe", func(c fiber.Ctx) error {
			req := ping.NewRequest()

			if err := c.Bind().Query(req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(
					ping.Response{Status: ping.StatusError, Message: err.Error()},
				)
			}

			if err := req.Validate(); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(
					ping.Response{Status: ping.StatusError, Message: err.Error()},
				)
			}

			resp, err := ping.Probe(c.Context(), req)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(resp)
			}

			return c.JSON(resp)
		})

		return app.Listen(args[0], fiber.ListenConfig{
			DisableStartupMessage: true,
			BeforeServeFunc: func(app *fiber.App) error {
				slog.Info("starting agent", slog.String("address", args[0]))
				return nil
			},
		})
	},
}
