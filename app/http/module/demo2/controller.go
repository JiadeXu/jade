package demo2

//func FooControllerHandler(c *gin.Context) error {
//	finish := make(chan struct{}, 1)
//	panicChan := make(chan interface{}, 1)
//
//	durationCtx, cancel := context.WithTimeout(c.BaseContext(), 1*time.Second)
//	defer cancel()
//
//	go func() {
//		defer func() {
//			if p := recover(); p != nil {
//				panicChan <- p
//			}
//		}()
//		// do real action
//		time.Sleep(10 * time.Second)
//		c.ISetOkStatus().IJson("ok")
//
//		finish <- struct{}{}
//	}()
//
//	select {
//	case p := <-panicChan:
//		c.WriterMux().Lock()
//		defer c.WriterMux().Unlock()
//		log.Println(p)
//		c.SetStatus(500).Json( "panic")
//	case <-finish:
//		fmt.Println("finish")
//	case <-durationCtx.Done():
//		c.WriterMux().Lock()
//		defer c.WriterMux().Unlock()
//		c.SetStatus(500).Json("time out")
//		c.SetHasTimeout()
//	}
//	return nil
//}
