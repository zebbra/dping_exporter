package cmd

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/spf13/cobra"
	"github.com/zebbra/dping_exporter/lib/ping"
)

func init() {
	rootCmd.AddCommand(agentCmd)
}

var agentCmd = &cobra.Command{
	Use:   "agent <listen-address>",
	Short: "Runs polling agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fiber.New()

		app.Get("/", func(c fiber.Ctx) error {
			return c.SendString("ready for service...")
		})

		app.Get("/-/healthy", healthcheck.New(healthcheck.Config{
			Probe: func(c fiber.Ctx) bool {
				return false
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

		return app.Listen(args[0])
	},
}
