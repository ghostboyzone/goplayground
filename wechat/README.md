### 微信网页版API

* usage

  ```
  import (
  	myApi "github.com/ghostboyzone/goplayground/wechat/api"
  )

  wechat := myApi.NewWechat()
  wechat.ShowQrCode()
  wechat.WaitForScan()
  wechat.GetContact()

  wechat.SendMsg("123", "filehelper")
  ```

  ​