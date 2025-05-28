文档
=============
TopChange.Net
只支持 IRR (伊朗里亚尔)

整体说明
=============
1. 发请pre-order请求, 返回一个token
2. 随后fe前端带着这个token, Get请求 (https://pg.toppayment.net/Merchant?Token=***  |   https://sandbox.toppayment.net?Token=***) 来打开收银台
3. 随后用户在收银台地址进行支付
4. merchant会收到对应的callback回调,收到后：merchant还需要发一个verify请求给psp来最终确认这一笔交易.

鉴权
==============
1. rsa privateKey私钥加密算签名, rsa publicKey公钥解密验证签名
2. 对请求参数算了一个sign签名, 随后作为json里一个字段一起发

回调地址
==============
api中可以指定回调地址, 所以是动态的


Comment
===============
1. 支持 deposit/withdraw
2. 所有接口都是 application/json 格式的
