已完成用户，视频，互动，社交模块

ID部分采用雪花ID

注册登录采用的是双token，用cookie存储

保护视频接口部分用了mysql的事务，保证上传，更新，删除视频时保存文件和数据库写入的一致性
删除视频以后将有关该视频缓存自动删除

上传视频，头像，封面做到可以正常接收文件，并保存到本地某个目录下

视频热门排行榜按照点赞数 x 100 + 评论数 x 20 来进行排序，并且通过zst存储videoid和hot_score与筛选前十热度视频

点赞数利用string类型缓存，点赞与取消点赞时通过UpdateRankScore更新排行榜的zset，并将变化点赞数的视频加入DirtyRedis后续进行异步写入mysql
点赞关系的删除与建立是直接写入mysql的

评论数利用string类型缓存，添加与删除时通过UpdateRankScore更新排行榜zset，并将变化评论数的视频加入DirtyRedis后续进行异步写入mysql
评论内容的添加与删除是直接写入mysql的
评论列表采用分页查询与string类型缓存(TTL30s)，在评论删除与添加后直接清除评论列表缓存

点赞和评论缓存预热采用 cache.TryWarmupLock(lockKey, 5*time.Second) 分布式锁+双重检查，防止并发预热导致重复 DB 查询


社交部分就是普通的curd

完成Docker部署，并且传到了docker hub上，就是pull下来总是启动失败(研究ing)
