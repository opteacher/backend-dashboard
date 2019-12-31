const fs = require("fs")
const _ = require("lodash")
const Path = require("path")
const router = require("koa-router")()

const scanPath = require("../utils/system").scanPath

// @block{userRoutes}:用户自定义路由
// @includes:path/../
// @includes:koa-router
// @includes:../utils/system.scanPath
console.log("API路由：")

// @steps{1}:扫描当前路径下除index.js以外的所有文件
let subPathAry = scanPath(__dirname, {ignores: ["index.js", "template.js"]})

// @steps{2}:根据各个文件相对路径，require之后开辟相应的路由
subPathAry.map(file => {
    let pthStat = Path.parse(file)
    let preRoutePath = `/${pthStat.dir.replace(/\\/g, "/")}`
    let refIdx = require(`./${file}`)
    let content = fs.readFileSync(Path.resolve(__dirname, file), "utf8")
    for(let i = content.indexOf("router."); i !== -1; i = content.indexOf("router.", i)) {
        i += "router.".length
        let bracket = content.indexOf("(", i)
        if(bracket === -1) {
            continue
        }
        let comma = content.indexOf(",", bracket)
        if(comma === -1) {
            continue
        }
        let subRoutePath = content.substring(bracket + 1, comma)
        subRoutePath = subRoutePath.substring(1, subRoutePath.length - 1)
        let routePath = ""
        if(subRoutePath !== "/") {
            routePath = preRoutePath + subRoutePath
        } else {
            routePath = preRoutePath
        }
        let method = content.substring(i, bracket).toLocaleUpperCase()
        console.log(`${method}${method.length > 3 ? "" : "\t"}\t${routePath}`)
    }
    router.use(preRoutePath, refIdx.routes(), refIdx.allowedMethods())
})

// @steps{3}:将index.html视图文件作为根路由/的渲染页面
router.get("/", async ctx => await ctx.render("index"))

module.exports = router