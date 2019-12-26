const puppeteer = require("puppeteer")
const config = require("./configs/puppeteer.json")

class Puppeteer {
    
    openChrome(show) {
        return puppeteer.launch({
			executablePath: config.chromium,
			headless: show
        })
    }

    async newPage(browser) {
        const page = await browser.newPage()
        page.setDefaultNavigationTimeout(config.timeout)
        return Promise.resolve(() => page)
    }

    close(browser) {
        return browser.close()
    }

    goto(page, url) {
        console.log(`跳转至${url}……`)
		return page.goto(url, {
			waitUntil: "networkidle2" // 等待网络状态为空闲的时候才继续执行
		})
    }

    wait(page, time) {
        console.log(`等待${time}ms`)
		return page.waitFor(time)
    }

    waitForElement(page, selector) {
        return page.waitForXPath(selector)
    }

    classStartsWith(page, prefix) {
        return page.waitForXPath(`//[starts-with(@class,'${prefix}')]`)
    }

    classContains(page, text) {
        return page.$x(`//[contains(@class,'${text}')]`)
    }

    querySelector(element, selector) {
        return element.$(selector)
    }

    querySelectorAll(element, selector) {
        return element.$$(selector)
    }

    async queryAndAdjustSelector(element, selector, procFunc) {
        let proceds = []
        await element.$eval(selector, ele => {
            proceds.push(procFunc(ele))
        })
        return Promise.resolve(() => proceds)
    }

    tap(element) {
        return element.tap()
    }

    type(element, text) {
        return element.type(text)
    }

    click(element) {
        return element.click()
    }
}

export default Puppeteer