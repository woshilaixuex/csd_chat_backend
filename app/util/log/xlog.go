package xlog

import (
	"log/slog"
	"os"
)

/*
 * @Author: Elyr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 基于slog基础上的自定义log
 * @Date: 2025-03-23 19:53
 */

func init() {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
