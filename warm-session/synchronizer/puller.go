package synchronizer

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func RunPullerLoop(
	ctx context.Context,
	sessionChan chan Session,
	completedChan chan string,
) {
	go func() {
		for {
			select {
			case <-ctx.Done(): // shut down loop
				close(completedChan) // the writer is supposed to close channels
				return
			case ses := <-sessionChan:
				go func() {
					fmt.Printf("puller: asked to pull image %s with policy %s and timeout %v\n",
						ses.image, ses.pullPolicy, ses.timeout)
					//TODO: where is pullPolicy used to make decisions? here?
					ctx2, _ := context.WithTimeout(ctx, ses.timeout)
					mustPull := !cri.hasImage(ses.image)
					pullStart := time.Now()
					if mustPull {
						fmt.Printf("puller: image not found, pulling %s\n", ses.image)
						cri.pullImage(ses.image, ctx2)
					}
					// update fields
					select {
					case <-ctx.Done(): // shutting down
						ses.isComplete = false
						ses.isTimedOut = false
						ses.err = fmt.Errorf("puller: shutting down")
						fmt.Println(ses.err.Error())
					case <-ctx2.Done():
						ses.isComplete = false
						ses.isTimedOut = true
						ses.err = fmt.Errorf("puller: async pull exceeded timeout of %v for image %s", ses.timeout, ses.image)
						fmt.Println(ses.err.Error())
					default:
						ses.isComplete = true
						ses.isTimedOut = false
						ses.err = nil
						if mustPull {
							fmt.Printf("puller: pull completed in %v for image %s\n", time.Since(pullStart), ses.image)
						} else {
							fmt.Printf("puller: image already present for %s\n", ses.image)
						}
					}
					close(ses.done)            // signal done
					completedChan <- ses.image //TODO: this could error if already closed above... that's ok
				}()
			}
		}

	}()
}

type criMock interface {
	hasImage(image string) bool
	pullImage(image string, ctx context.Context)
}

var cri criMock = criMockImpl{}

type criMockImpl struct{}

func (c criMockImpl) hasImage(image string) bool { return false }
func (c criMockImpl) pullImage(image string, ctx context.Context) {
	min := 3000
	max := 7000
	rand.Seed(time.Now().UnixNano()) // without seed, same sequence always returned
	dur := time.Duration(min + rand.Intn(max-min))

	fmt.Printf("criMock: starting to pull image %s\n", image)
	select {
	case <-time.After(dur * time.Millisecond):
		fmt.Printf("criMock: finshed pulling image %s\n", image)
		return
	case <-ctx.Done():
		fmt.Printf("criMock: context cancelled\n")
		return
	}
}
