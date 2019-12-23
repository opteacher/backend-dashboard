<template>
<el-form :inline="true" ref="load-dao-group" :model="selDaoGroup">
    <el-form-item label="语言">
        <el-select v-model="selLang" @change="chgLangOrCate">
            <el-option v-for="lang in supportLangs" :key="lang" :label="lang" :value="lang"/>
        </el-select>
    </el-form-item>
    <el-form-item label="类别">
        <el-select v-model="selCategory" @change="chgLangOrCate">
            <el-option v-for="catgry in supportCategories" :key="catgry" :label="catgry" :value="catgry"/>
        </el-select>
    </el-form-item>
    <el-form-item class="float-right" label="导入的DAO组">
        <el-select v-model="selGrpName" @change="chgDaoGroup">
            <el-option v-for="group in filterTmpDaoGrps" :key="group.name" :label="group.name" :value="group.name"/>
        </el-select>
    </el-form-item>
</el-form>
</template>

<script>
import _ from "lodash"
import backend from "../backend"

export default {
    data() {
        return {
            tempDaoGroups: [],
            filterTmpDaoGrps: [],
            supportLangs: [],
            supportCategories: [],
            selLang: "*",
            selCategory: "*",
            selDaoGroup: null,
            selGrpName: ""
        }
    },
    async created() {
        let res = await backend.qryAllTempDaoGroups()
        if (typeof res === "string") {
            this.$message.error(`查询所有模板DAO组时发生错误：${res}`)
        } else {
            this.tempDaoGroups = res.groups
            this.filterTmpDaoGrps = _.cloneDeep(this.tempDaoGroups)
            let sptLangs = ["*"]
            let sptCatgrs = ["*"]
            for (let group of this.tempDaoGroups) {
                sptLangs.push(group.language)
                for (let category of group.categories) {
                    sptCatgrs.push(category)
                }
            }
            this.supportLangs = _.uniq(sptLangs)
            this.supportCategories = _.uniq(sptCatgrs)
        }
    },
    methods: {
        chgLangOrCate() {
            this.filterTmpDaoGrps = _.cloneDeep(this.tempDaoGroups)
            if (this.selLang !== "*") {
                this.filterTmpDaoGrps = _.filter(this.filterTmpDaoGrps, grp => grp.language === this.selLang)
            }
            if (this.selCategory !== "*") {
                this.filterTmpDaoGrps = _.filter(this.filterTmpDaoGrps, grp => grp.categories.includes(this.selCategory))
            }
        },
        chgDaoGroup(selected) {
            this.selDaoGroup = this.tempDaoGroups.find(grp => grp.name === selected)
        }
    }
}
</script>