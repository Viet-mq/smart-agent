package main

//func main() {
//	cst := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ts := r.ServerTraceState()
//		fmt.Printf("Request was received at: %s.\n\n", ts.ReceivedAt)
//		fmt.Printf("TLS handshaking started after: %s.\nTime for TLS handshake: %s\n\n",
//			ts.TLSHandshakeStartAt.Sub(ts.ReceivedAt),
//			ts.TLSHandshakeEndAt.Sub(ts.TLSHandshakeStartAt))
//		fmt.Printf("Header processing started after: %s.\nFinished after: %s.\nTotal time to parse headers: %s.\n\n",
//			ts.ParsedHdrsStartAt.Sub(ts.ReceivedAt),
//			ts.ParsedHdrsEndAt.Sub(ts.ReceivedAt),
//			ts.ParsedHdrsEndAt.Sub(ts.ParsedHdrsStartAt))
//		fmt.Printf("Reading the first byte started after: %s.\nTime to read first byte: %s.\n\n",
//			ts.FirstByteStartAt.Sub(ts.ReceivedAt),
//			ts.FirstByteEndAt.Sub(ts.FirstByteStartAt))
//		fmt.Printf("Composing the request before passing it into this handler took: %s.\n",
//			ts.ParsedRequestAt.Sub(ts.ReceivedAt))
//	}))
//	defer cst.Close()
//
//	req, _ := http.NewRequest("GET", cst.URL, nil)
//	res, err := cst.Client().Do(req)
//	if err != nil {
//		panic(err)
//	}
//	_, _ = ioutil.ReadAll(res.Body)
//	_ = res.Body.Close()
//}
