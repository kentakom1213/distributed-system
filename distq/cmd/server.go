package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	queueserver "github.com/kentakom1213/distributed-system/distq/internal/server"
	"github.com/kentakom1213/distributed-system/distq/internal/storage"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the queue server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
		defer stop()
		
		st, err := storage.Open(ctx, dbPath)
		if err != nil {
			return err
		}
		defer st.Close()

		handler := queueserver.NewHandler(st)

		srv := &http.Server{
			Addr:              serverAddr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		}

		errCh := make(chan error, 1)

		go func() {
			slog.Info("queue server started", "addr", serverAddr, "db", dbPath)
			errCh <- srv.ListenAndServe()
		}()

		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			slog.Info("queue server shutting down")
			return srv.Shutdown(shutdownCtx)
		case err := <-errCh:
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		}
	},
}

var (
	serverAddr string
	dbPath     string
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&serverAddr, "addr", ":8080", "server address")
	serverCmd.Flags().StringVar(&dbPath, "db", "data/distq.db", "sqlite database path")
}
