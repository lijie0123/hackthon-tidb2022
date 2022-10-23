<h1 align="center">Open Source Crypto Insight!</h1>

![demo](https://raw.githubusercontent.com/lenny-mo/PictureUploadFolder/main/Screen%20Shot%202022-10-23%20at%2011.48.10.png)

![](https://raw.githubusercontent.com/lenny-mo/PictureUploadFolder/main/Screen%20Shot%202022-10-23%20at%2012.05.58.png)

![](https://raw.githubusercontent.com/lenny-mo/PictureUploadFolder/main/Screen%20Shot%202022-10-23%20at%2011.47.20.png)

<p align="center">
<b>Open Source Crypto Insight</b> analyzes billions of data of crypto and gives a professional perspective based on these data. 
</p>

<p align="center">
It also provides a bunch of utensils which can show different perspectives for different users, which will be detailed next section.
</p>

## Magic Dashboard👁️

经过筛选的数据，分不同目的，分区块展示。

### 1.0 Professional Vision

The dashboard, which has been screened by a professional team, offers a subscription service and sends an alert to the user's mailbox when the alert threshold set by the user is reached.

经过专业团队人员筛选过的 Dashboard,提供订阅服务，在达到用户设置的警戒线的时候向用户邮箱发出警报。

### 1.1 Real-time Chart

Real-time display of market prices and trading volumes for cryptocurrencies.

实时展示加密货币的市场价格和交易量。

### 1.2 Market Data for Beginner

Show foundamental details related to the crypto which is selected by you, which includes definition, statistics analysis, and market information.

展示加密货币的基本信息，包括但是不限于加密货币的官方定义，数值统计分析和市场对该加密货币的舆情。

### 1.3 Community for socialization

This section includes formal and informal information in Reddit, Twitter, and Meta(also known as Facebook).

这个部分展示加密货币在国外主流媒体平台上的评价信息，包括但是不限于 Reddit, Twitter, Meta Groups.

### 1.4 Connect to OSS Insight effortlessly

In this section we show the open source projects for cryptocurrencies, including the number of stars and forks of this project on Github. Based on this publicly available data, we use a set of algorithms to calculate the reliability of cryptocurrencies.

在这个板块，我们展示了加密货币的开源项目，包括这个项目在 Github 上的 Star 数量，Fork 数量。根据这些可以获取到的公开数据，我们通过一套算法，计算出该加密货币的可靠性。

## API for using

后续会陆续开放各种 Blockchain 上获取数据的 API。

- just a tiny example:

by running the code right below,

```sql
curl http://localhost:1323/api/btc/blocks/00000000000000000006ba2a50ae990822cf8fbd4b22398b914703c0275e6754

```

and then, you can get data like this:

```
{"bits":"1707e772","coinbase_param":"03d3970b1b4d696e656420627920416e74506f6f6c383439b201af030b8a11abfabe6d6d314c2d5c03754fcee7344d9c3a7f6945f2a34a1d62b9d2aa3c010adb5757634002000000000000000000ef87bb00000000000000","hash":"00000000000000000006ba2a50ae990822cf8fbd4b22398b914703c0275e6754","merkle_root":"7478debd909563ad3a9c62401b7ba11436338bd779e5d1affce2e756f7fa27ec","nonce":"9269a854","number":759763,"size":305490,"stripped_size":175829,"timestamp":"2022-10-22 02:38:17","transaction_count":547,"version":705691648,"weight":832977}
```

## Development

- you can click [here](http://ec2-35-77-75-24.ap-northeast-1.compute.amazonaws.com/coins/ethereum) to access to our website.

## Sponsors

- tiDB
