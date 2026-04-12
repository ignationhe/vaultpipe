package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
/envfile"
)
atchCmd :=ootCmd.AddCommand(watchCmdvar interval time.Duration
	var quiet bool

	cmd := &cobra.Command{
		Use:   "watch <file>",
		Short: "Watch a .env file and print changes as they occur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatch(args[0], interval, quiet)
		},
	}

	cmd.Flags().DurationVarP(&interval, "interval", "i", 2*time.Second, "polling interval")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress unchanged notifications")
	return cmd
}

func runWatch(path string, interval time.Duration, quiet bool) error {
	opts := envfile.DefaultWatchOptions()
	opts.Interval = interval
	opts.OnChange = func(p string, kv map[string]string) {
		fmt.Fprintf(os.Stdout, "[changed] %s — %d keys\n", p, len(kv))
		for k, v := range kv {
			fmt.Fprintf(os.Stdout, "  %s=%s\n", k, v)
		}
	}
	opts.OnError = func(p string, err error) {
		fmt.Fprintf(os.Stderr, "[error] %s: %v\n", p, err)
	}

	if !quiet {
		fmt.Fprintf(os.Stdout, "Watching %s every %s (Ctrl+C to stop)\n", path, interval)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := envfile.Watch(ctx, path, opts)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return nil
	}
	return err
}
