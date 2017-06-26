* api目录

    * 开放接口：[`api/jubi_public.go`](/bitcoin/api/jubi_public.go)

    * 需要进行权限验证接口：[`api/jubi.go`](/bitcoin/api/jubi.go)

* demo（需要先修改 [`api/const.go`](/bitcoin/api/const.go)里的key）

    * 拉取数据并保存到本地数据库文件：[`task_grabdata.go`](/bitcoin/task_grabdata.go)

    * 计算区间综合涨幅

        * 从开始日期到结束日期：[`task_rangeprofit.go`](/bitcoin/task_rangeprofit.go) 
          ```
          go run task_rangeprofit.go -start "2017-06-26 00:00" -end "2017-06-27 00:00"
          ```

        * 从开始日期到实时：[`task_realtimeprofit.go`](/bitcoin/task_realtimeprofit.go) 
          ```
          go run task_realtimeprofit.go -start "2017-06-26 00:00"
          ```
    
    * 推荐交易相关

        * 推荐买入：[`task_suggest.go`](/bitcoin/task_suggest.go) 
          ```
          go run task_suggest.go
          ```
        
        * 展示涨跌幅较大的币：[`task_autobuy.go`](/bitcoin/task_autobuy.go) 
          ```
          go run task_autobuy.go
          ```

    * 个人相关

        * 查询历史所有挂单：[`task_checkorders.go`](/bitcoin/task_checkorders.go) 
          ```
          go run task_checkorders.go
          ```

        * 查询实时收益情况：[`task_checkprofit.go`](/bitcoin/task_checkprofit.go) 
          (需要修改代码)
          ```
          go run task_checkprofit.go
          ```
