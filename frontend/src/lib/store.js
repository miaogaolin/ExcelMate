import { writable } from 'svelte/store';
import { ReadExcel, DialogError, ReadConfig, ReadSettings } from "../../wailsjs/go/main/App";

export const allColor = [
    "rgb(99 102 241)",
    "rgb(180 83 9)",
    "rgb(134 239 172)",
    "rgb(15 118 110)",
    "rgb(29 78 216)",
    "rgb(216 180 254)",
    "rgb(112 26 117)",
    "rgb(157 23 77)",
    "rgb(217 119 6)",
    "rgb(250 204 21)",
    "rgb(239 68 68)",
    "#f08080",
    "#ff00ff",
    "rgb(103 232 249)",
];

export const configData = writable({
    current: "default",
    files: [
        {
            id: "default",
            name: "Untitled",
            list: [{

                condition: "",
                template: "{{.A}}",
                color: allColor[0]
            }]
        }
    ]
});


export function initExcelData() {
    ReadExcel("").then(r => {
        if (r && r.length > 0) {
            let tmp = [];
            r.forEach((v) => {
                tmp.push({
                    data: v,
                });
            });
            excelData.set(tmp);
        }
    }).catch(e => {
        DialogError(e);
    })
}

export function initConfigData() {
    ReadConfig("")
        .then(r => {
            if (r != "") {
                console.log("init config", JSON.parse(r));
                configData.set(JSON.parse(r));
            }
        })
        .catch((e) => {
            DialogError(e);
        });
}

export function initSettingsData() {
    ReadSettings().then((r) => {
        settingsData.set(r);
    });
}

export function init() {
    // 数据配置
    initConfigData();
    // 初始化excel数据
    initExcelData();
    // 应用设置
    initSettingsData();
}


export const excelData = writable([]);
export const settingsData = writable({
    config_file: "",
    excel_file: "",
});

// 模板匹配数据的结果
export const outputData = writable({
    config_index: -1, // 仅展示对应条件的结果, -1 代表所有
    list: []
});


