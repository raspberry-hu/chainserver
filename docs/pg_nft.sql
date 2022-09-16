CREATE USER mintklubdev WITH PASSWORD 'mintklub.com';        -- 创建 mintklubdev 用户， 密码 mintklub.com
CREATE DATABASE mintklubdb OWNER mintklubdev;                -- 创建 mintklubdb 数据库 所有权限都赋予 mintklubdev 用户
\c mintklubdb                                              -- mintklubdb 数据库
CREATE SCHEMA mintklub;                                    -- 创建 mintklub 模式
GRANT ALL ON schema mintklub to mintklubdev; -- 将模式权限赋值给 mintklubdev
ALTER USER mintklubdev SET search_path to mintklub;   -- 设置search_path
show search_path;   -- 查看

GRANT ALL ON schema mintklub to mintklubdev; -- 将模式权限赋值给 mintklubdev

CREATE TABLE mintklub.t_market_list (
      ID serial PRIMARY KEY, --自增长 int 类型，主键
      nft_id INT NOT NUll, -- nft id
      creater VARCHAR (256) NOT NULL, -- 挂单者,
      token_id VARCHAR (256) NOT NULL, -- 挂单nft编号,
      market_type INT NOT NULL DEFAULT 0, -- 挂单类型：2 限价购买；1 拍卖,
      starting_price VARCHAR (256) NOT NULL, -- 起拍价,
      end_time INT NOT NULL, -- 拍卖结束时间,
      buyer VARCHAR (100) DEFAULT '', -- 成功购买者,
      reward INT NOT NULL DEFAULT 0, -- 拍卖参与者分红,
      tx_hash VARCHAR (256) DEFAULT '', -- 交易hash,
      cancel_hash VARCHAR (256) DEFAULT '', -- 取消挂单hash,
      deal_hash VARCHAR (256) DEFAULT '', -- deal hash,
      create_time INT DEFAULT 0, -- 挂单时间,
      status INT DEFAULT 0, -- 状态：0 已提交 1 已上链 2 失败 3 取消挂单 4 挂单成交,
      chain_name VARCHAR (256) NOT NULL, -- 区块链名称
      currency_name VARCHAR (256) NOT NUll, --币种
      lazy INT DEFAULT 0, -- 0 收取gas费、 1不收
      donation INT DEFAULT 0, -- 0 不捐赠、1捐赠
      donation_user_id VARCHAR (256) DEFAULT '' -- 捐赠者用户id
);

CREATE TABLE mintklub.t_nft (
      ID serial PRIMARY KEY, --自增长 int 类型
      sn VARCHAR (256) NOT NULL, -- 流水号
      nft_name VARCHAR (256) NOT NULL, -- 名称
      nft_desc TEXT NOT NULL, -- 描述
      rights_rules TEXT DEFAULT '', -- 权益规则
      token_id VARCHAR (256) DEFAULT '', -- nft编号
      meta_data_uri TEXT DEFAULT '', -- meta_data_uri 原数据地址
      explore_uri TEXT DEFAULT '', -- 音视频的数据地址
      tx_hash VARCHAR (256) DEFAULT '', -- 铸造nft的hash
      transfer_hash VARCHAR (256) DEFAULT '', -- 直接转移 nft 的hash
      creater VARCHAR (256) NOT NULL, -- nft创建者
      block_number INT DEFAULT 0, -- 创建区块高度
      create_time INT DEFAULT 0, -- 创建时间戳
      media_uri TEXT NOT NULL, -- nft 图片地址
      create_tax INT NOT NULL, -- 铸造税
      owner INT NOT NULL, -- nft拥有者
      status INT DEFAULT 0, -- 状态：0 已提交 1 已上链 2 失败
      market_type INT DEFAULT 0, -- 挂单状态：0 未挂单 1 拍卖 2 限价购买
      approved INT DEFAULT 0, -- 审核状态 0 未审核 1 审核过
      chain_name VARCHAR (256) NOT NULL, -- 区块链名称
      currency_name VARCHAR (256) NOT NUll, --币种目前支持ETH\AVAX\BNB
      lazy INT DEFAULT 0, -- 0 收取gas费、 1不收
      media_ipfs_uri TEXT DEFAULT '', -- ipfs_uri ipfs 图片地址
      collection_id INT DEFAULT 0, -- 集合id
      categories_id INT DEFAULT 0 -- 类别id
);

CREATE TABLE mintklub.t_order (
      ID serial PRIMARY KEY, --自增长 int 类型，主键
      nft_id INT NOT NULL, -- 对应的nft id
      market_id INT NOT NULL, -- 对应的挂单id
      market_type INT DEFAULT 0, -- 订单类型 竞拍 1 限价 2
      token_id VARCHAR (256) NOT NULL, -- tokenid
      create_time INT NOT NULL, -- 下单时间
      onchain_time INT DEFAULT 0, -- 上链时间
      tx_hash VARCHAR (256) DEFAULT '', -- 交易hash
      buyer VARCHAR (256) NOT NULL, -- 买家
      seller VARCHAR (256) NOT NULL, -- 卖家
      price VARCHAR (256) NOT NULL, -- 价格
      chain_name VARCHAR (256) NOT NULL, -- 区块链名称
      currency_name VARCHAR (256) NOT NUll, --币种目前支持ETH\AVAX\BNB
      lazy INT DEFAULT 0, -- 0 收取gas费、 1不收
      status INT DEFAULT 0, -- 订单状态：0 已提交 1 已上链 2 失败
      donation INT DEFAULT 0 -- 0 不捐赠、1捐赠
);

CREATE TABLE mintklub.t_wallet_info (
  ID serial PRIMARY KEY, --自增长 int 类型 主键
  wallet_addr VARCHAR (256) NOT NULL DEFAULT '', -- 钱包地址
  user_name VARCHAR (256) NOT NULL DEFAULT '', -- 用户名
  user_desc text DEFAULT '', -- 用户描述
  image_url VARCHAR (256) NOT NULL DEFAULT '', -- 头像地址
  banner_url VARCHAR (256) NOT NULL DEFAULT '', -- 横幅地址
  email_addr  VARCHAR (256) NOT NULL DEFAULT '', -- email
  create_time TIMESTAMP (6) WITHOUT TIME ZONE DEFAULT (now()), -- 创建时间
  update_time TIMESTAMP (6) WITHOUT TIME ZONE DEFAULT (now()) -- 修改时间
);

CREATE TABLE mintklub.t_favorites (
   ID serial PRIMARY KEY, --自增长 int 类型 主键
   user_id VARCHAR (256) NOT NULL, -- 用户id
   nft_id  int , -- nft 主键id
   create_time TIMESTAMP (6) WITHOUT TIME ZONE DEFAULT (now()) -- 创建时间
);

CREATE TABLE mintklub.t_collection_info (
   ID serial PRIMARY KEY, --自增长 int 类型 主键
   user_id VARCHAR (256) NOT NULL, -- 用户id
   collection_name VARCHAR (256),  -- 集合名称
   chain_name VARCHAR (256) NOT NULL, -- 区块链名称
   currency_name VARCHAR (256) NOT NUll, --币种目前支持ETH\AVAX\BNB
   logo_image VARCHAR (256) NOT NUll, -- 形象标识
   featured_image_url VARCHAR (256) NOT NUll, -- 特色图片3
   banner_image_url VARCHAR (256) NOT NUll, -- 横幅图片
   collection_desc VARCHAR (256) NOT NUll, -- 集合描述
   create_tax INT DEFAULT 0, -- 铸造税 //1
   category VARCHAR (256) NOT NULL, --类别名称
   items INT DEFAULT 0, -- 资产总量
   favorites INT DEFAULT 0, -- 收藏资产总量
   amount FLOAT DEFAULT 0, -- 成交总额
   tx_hash VARCHAR (256) DEFAULT '', -- 交易hash,
   status INT DEFAULT 0 -- 状态：0 已提交 1 已上链 2 失败 3 取消挂单 4 挂单成交,
);

CREATE TABLE mintklub.t_categories_info (
     ID serial PRIMARY KEY, --自增长 int 类型 主键
     categories_name VARCHAR (256), -- 类别名称
     categories_desc VARCHAR (256), -- 标签描述
     collection_id INT, -- 集合id
     user_id VARCHAR (256) NOT NULL -- 用户id
);

CREATE TABLE mintklub.t_organization_info (
   ID serial PRIMARY KEY, --自增长 int 类型 主键
   org_name VARCHAR (256), -- 组织名
   image_url VARCHAR (256), -- 头像地址
   user_id VARCHAR (256) -- 用户id
);

-- 捐赠记录
CREATE TABLE mintklub.t_donation_record (
    ID serial PRIMARY KEY, --自增长 int 类型 主键
    nft_id INT NOT NUll, -- nft id
    nft_name VARCHAR (256) NOT NULL, -- 名称
    nft_desc TEXT NOT NULL, -- 描述
    market_id INT NOT NUll, -- nft id
    token_id VARCHAR (256) NOT NULL, -- 挂单nft编号,
    buyer VARCHAR (256) NOT NULL, -- 买家
    seller VARCHAR (256) NOT NULL, -- 卖家 捐赠者
    donation_amount VARCHAR (256) NOT NULL, -- 捐赠金额
    create_time INT DEFAULT 0 -- 挂单时间
);

CREATE TABLE mintklub.t_permissions (
    ID serial PRIMARY KEY, --自增长 int 类型 主键
    user_id VARCHAR (256) -- 用户id
);


-- CREATE TABLE "mintklub"."t_user_info" (
--    ID serial PRIMARY KEY, --自增长 int 类型 主键
--   "user_name" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
--   "user_pass" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
--   "user_desc" text COLLATE "pg_catalog"."default" DEFAULT ''::text,
--   "image_url" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
--   "banner_url" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
--   "email_addr" varchar(256) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
--   "create_time" timestamp(6) DEFAULT now(),
--   "update_time" timestamp(6) DEFAULT now()
-- );

CREATE TABLE "mintklub"."t_user_info" (
    ID serial PRIMARY KEY, --自增长 int 类型 主键
    "user_name" varchar(256) NOT NULL,
    "user_pass" varchar(256) NOT NULL,
    "user_desc" text COLLATE "pg_catalog"."default" DEFAULT ''::text,
    "image_url" varchar(256) NOT NULL,
    "banner_url" varchar(256) NOT NULL,
    "email_addr" varchar(256) NOT NULL,
    "create_time" timestamp(6) DEFAULT now(),
    "update_time" timestamp(6) DEFAULT now()
);