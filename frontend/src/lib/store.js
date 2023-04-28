import { writable } from 'svelte/store';
import { SaveConfig, ReadExcel, DialogError, ReadConfig, ReadSettings } from "../../wailsjs/go/main/App";

export const configData = writable({
    current: "default",
    files: [
        {
            id: "default",
            name: "Untitled",
            list: [{

                condition: "",
                template: "{{.A}}",
                color: "border-s-indigo-500"
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
export const outputData = writable([]);

