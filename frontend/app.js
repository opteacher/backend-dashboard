const path = require("path");
const Koa = require("koa");
// 用于旧版本模块的generator向async/await转换
const convert = require("koa-convert");
const bodyparser = require("koa-bodyparser");
const json = require("koa-json");
const logger = require("koa-logger");
const statc = require("koa-static");
const views = require("koa-views");
const cors = require("koa2-cors");

const app = new Koa();

// 跨域配置
app.use(cors());

// 路径解析
app.use(bodyparser());

// json解析
app.use(json());

// 日志输出
app.use(logger());

// 指定静态目录
app.use(statc(path.join(__dirname, "public")));

// 指定页面目录
app.use(views(path.resolve("./"), {extension: "html"}));

// 路径分配
app.use(async ctx => await ctx.render("index"));

// 监听
app.listen(process.env.PORT || 4000);