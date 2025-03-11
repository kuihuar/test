package redpackage

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Listen() {
	// 创建监听器
	listener, _ := net.Listen("tcp", ":5577")

	defer listener.Close()

	// 使用WaitGroup跟踪活跃连接
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Printf("Received signal: %v. Shutting down...\n", sig)

		listener.Close()
		cancel()
		time.AfterFunc(time.Second*10, func() {
			fmt.Println("Timeout reached. Forcing shutdown...")
			os.Exit(1)
		})
	}()
	fmt.Println("Server started. Press Ctrl+C to exit.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("Stop accepting new connections.")
				break
			default:
				fmt.Println("Accept error: %v\n", err)
			}
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(ctx, conn)
		}()
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Received: %s\n", string(buf[:n]))
	}

	conn.Write([]byte("Hello World"))
}
