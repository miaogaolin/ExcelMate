<script>
    import { ValidateExpr } from "../wailsjs/go/main/App";
    import {
        configData,
        allColor,
        outputData,
        conditionError,
        templateError,
    } from "./lib/store.js";
    import { getUuid, selectText } from "./lib/utils.js";
    import Textarea from "./component/Textarea.svelte";

    // 当前选择的配置文件
    let currentConfig = { list: [] };
    // 所有配置
    let allConfig = { current: "", files: [] };
    // 已使用的颜色
    let usedColor = [];
    // 是否展示列表
    let isList = false;
    // 输出结果展示的索引
    let outputDataIndex = -1;
    $: {
        outputData.update((r) => {
            r.config_index = outputDataIndex;
            return r;
        });
    }

    function updateConfig(index, conditionEle, templateEle) {
        let condition, template;

        if (conditionEle) {
            condition = conditionEle.target["value"].trim();
        }

        if (templateEle) {
            template = templateEle.target["value"];
        }

        if (
            (condition != undefined &&
                currentConfig.list[index].condition == condition) ||
            (template != undefined &&
                currentConfig.list[index].template == template)
        ) {
            return;
        }

        configData.saveUpdate((all) => {
            let current = getCurrentConfig(all);
            let list = current.list;
            if (index < list.length) {
                let d = list[index];
                if (condition != undefined) {
                    d.condition = condition;
                }

                if (template != undefined) {
                    d.template = template;
                }
                list[index] = d;
                current.list = list;
            }
            return all;
        });
    }

    function addConditionAndTemplate() {
        configData.saveUpdate((all) => {
            let color = getRandomColor();
            for (let i = 0; i < allColor.length; i++) {
                if (!usedColor.includes(allColor[i])) {
                    color = allColor[i];
                    break;
                }
            }

            let current = getCurrentConfig(all);
            current.list.push({
                condition: "",
                template: "",
                color: color,
            });
            return all;
        });
    }

    function deleteConditionAndTemplate(index) {
        configData.saveUpdate((all) => {
            let current = getCurrentConfig(all);
            let list = current.list;
            if (index < list.length) {
                list.splice(index, 1);
                current.list = list;
            }
            return all;
        });
    }

    function getRandomColor() {
        // 随机生成 RGB 值
        const r = Math.floor(Math.random() * 256);
        const g = Math.floor(Math.random() * 256);
        const b = Math.floor(Math.random() * 256);

        // 将 RGB 转换为十六进制颜色值，并返回结果
        return (
            "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)
        );
    }

    function getCurrentConfig(data) {
        for (let i = 0; i < data.files.length; i++) {
            if (data.current == data.files[i].id) {
                return data.files[i];
            }
        }
        return {
            list: [],
        };
    }

    function newConfig() {
        let id = getUuid();
        let newConfig = {
            id: id,
            name: "Untitled",
            list: [
                {
                    condition: "",
                    template: "{{.A}}",
                    color: allColor[0],
                },
            ],
        };

        configData.saveUpdate((all) => {
            all.files.push(newConfig);
            all.current = id;
            return all;
        });

        let configName = document.getElementById("ConfigName");
        if (configName) {
            configName.focus();
        }
    }

    function updateConfigName(e) {
        let text = e.target.innerText;
        if (text == "") {
            return;
        }

        configData.saveUpdate((all) => {
            let current = getCurrentConfig(all);
            current.name = text;
            return all;
        });
    }

    function validateCondition(el, index) {
        let condition = el.target["value"].trim();
        if (condition == "") {
            updateConfig(index, el, undefined);
            return;
        }

        // 单纯验证语法
        ValidateExpr(condition)
            .then(() => {
                updateConfig(index, el, undefined);
                $conditionError[index] = "";
            })
            .catch((e) => {
                $conditionError[index] = e;
            });
    }

    function exchangePosition(i, j) {
        configData.saveUpdate((all) => {
            let current = getCurrentConfig(all);
            if (i < 0 || j >= current.list.length) {
                return all;
            }

            let tmp = current.list[i];
            current.list[i] = current.list[j];
            current.list[j] = tmp;
            return all;
        });
    }

    configData.subscribe((v) => {
        if (v.current == "") {
            return;
        }
        console.log("configData.subscribe", v);
        allConfig = v;
        currentConfig = getCurrentConfig(v);
        let tmp = [];
        currentConfig.list.forEach((r) => {
            tmp.push(r.color);
        });
        usedColor = tmp;
    });

    // 多配置文件
    function selectConfig(id) {
        configData.saveUpdate((all) => {
            all.current = id;
            return all;
        });
        isList = false;
    }

    function deleteConfig(id) {
        configData.saveUpdate((all) => {
            let newIndex = 0;
            for (let i = 0; i < all.files.length; i++) {
                let row = all.files[i];
                if (row.id == id) {
                    newIndex = i;
                    all.files.splice(i, 1);
                    break;
                }
            }
            if (newIndex >= all.files.length) {
                newIndex--;
            }

            if (id == all.current) {
                if (newIndex >= 0) {
                    all.current = all.files[newIndex].id;
                } else {
                    all.files = [];
                    all.current = "";
                }
            }
            return all;
        });
    }

    function onlyShowOutput(i) {
        if (outputDataIndex == i) {
            outputDataIndex = -1;
            return;
        }
        outputDataIndex = i;
    }
</script>

<div
    class="condition-container w-full h-1/2 bg-slate-100 overflow-y-auto border-s-slate-100 dark:border-s-gray-900 border-s dark:bg-gray-900"
>
    <!-- 编辑配置 -->
    <div class:hidden={isList}>
        <div
            class="bg-white dark:bg-gray-900 flex items-center sticky top-0 z-50"
        >
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            {#if allConfig.current != ""}
                <span
                    id="ConfigName"
                    on:click={(e) => selectText(e.target)}
                    on:focus={(e) => selectText(e.target)}
                    on:focusout={updateConfigName}
                    class="text-xs p-2 font-medium bg-slate-100 dark:text-gray-400 dark:bg-gray-900 dark:border-gray-900 dark:hover:border-gray-700 dark:focus:border-gray-700 focus:outline-none border-2 border-slate-100 hover:border-slate-300 focus:border-slate-300"
                    contenteditable="true"
                >
                    {currentConfig.name}
                </span>
            {/if}
            <div class="ml-2 py-2">
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <svg
                    on:click={newConfig}
                    class="hover:cursor-pointer"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    ><path
                        d="M888.494817 313.882803l-198.019982-198.019982c-7.992021-7.992021-20.957311-7.992021-28.949332 0s-7.992021 20.947078 0 28.939099l163.084309 163.084309-215.794811 0L608.814999 42.686195c0-11.307533-9.15859-20.466124-20.466124-20.466124l-408.094512 0c-11.307533 0-20.466124 9.15859-20.466124 20.466124l0 938.62761c0 11.2973 9.15859 20.466124 20.466124 20.466124l693.76067 0c11.307533 0 20.466124-9.168824 20.466124-20.466124l0-652.961452C894.481158 322.92883 892.332215 317.720202 888.494817 313.882803zM853.54891 960.847681l-652.828422 0L200.720488 63.152319l367.162264 0 0 265.200034c0 11.307533 9.168824 20.466124 20.466124 20.466124l265.200034 0L853.54891 960.847681z"
                        fill="#8a8a8a"
                    /></svg
                >
            </div>
            <div class="ml-2">
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <svg
                    on:click={() => (isList = true)}
                    class="hover:cursor-pointer"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    ><path
                        d="M874.690416 495.535003c0 5.423523-2.159176 10.632151-5.996574 14.46955l-223.00912 223.00912c-4.001127 3.990894-9.240455 5.996574-14.46955 5.996574-5.239328 0-10.478655-1.995447-14.479783-5.996574-7.992021-7.992021-7.992021-20.947078 0-28.939099l188.083679-188.083679-604.773963 0c-11.2973 0-20.466124-9.168824-20.466124-20.466124 0-11.307533 9.168824-20.466124 20.466124-20.466124l604.753497 0-188.073446-188.073446c-7.992021-7.992021-7.992021-20.947078 0-28.949332 7.992021-7.992021 20.957311-7.992021 28.949332 0l223.019353 223.029586C872.53124 484.902852 874.690416 490.101247 874.690416 495.535003z"
                        fill="#8a8a8a"
                    /></svg
                >
            </div>
        </div>
        <div class="h-2" />
        {#each currentConfig.list as r, i}
            <div
                style="border-left-color:{r.color}"
                class="config-container border-l-4 overflow-hidden pr-6"
            >
                <div
                    class="float-left w-full {i === outputDataIndex
                        ? 'bg-gray-300 dark:bg-gray-600'
                        : 'bg-white dark:bg-gray-800'}"
                >
                    <Textarea
                        class="tracking-widest dark:font-light placeholder:italic placeholder:text-slate-400 block bg-transparent  dark:text-gray-200 dark:border-gray-700 w-full border border-slate-300 py-2 pl-2.5 pr-3 focus:outline-none focus:ring-1 sm:text-sm resize-none overflow-hidden"
                        value={r.condition}
                        rows="1"
                        placeholder="The conditions for filtering Excel."
                        on:change={(e) => validateCondition(e, i)}
                    />
                    {#if $conditionError[i] && $conditionError[i] != ""}
                        <div class="bg-red-200 text-xs p-2 my-2">
                            {$conditionError[i]}
                        </div>
                    {/if}

                    <Textarea
                        rows="3"
                        class="tracking-widest dark:font-light dark:text-gray-200 dark:border-gray-700 placeholder:italic placeholder:text-slate-400 block w-full border border-slate-300 py-2 pl-2.5 pr-3 focus:outline-none focus:ring-1 sm:text-sm resize-none overflow-hidden bg-transparent"
                        placeholder="The template content that matches the conditions successfully."
                        value={r.template}
                        on:change={(e) => updateConfig(i, undefined, e)}
                    />
                    {#if $templateError[i] && $templateError[i] != ""}
                        <div class="bg-red-200 text-xs p-2 my-2">
                            {$templateError[i]}
                        </div>
                    {/if}
                </div>
                <div class="float-right -mr-6 w-6 text-center">
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <svg
                        on:click={() => exchangePosition(i - 1, i)}
                        class:hidden={i == 0}
                        class="hover:cursor-pointer fill-gray-500 hover:fill-gray-700 dark:fill-gray-400 dark:hover:fill-white mx-auto mb-1"
                        viewBox="0 0 1024 1024"
                        version="1.1"
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
                        ><path
                            d="M919.194 500.382l-395.16-388.45c-6.649-6.536-17.419-6.536-24.068 0l-395.16 388.45a16.524 16.524 0 0 0 0 23.659c3.324 3.268 7.679 4.902 12.034 4.902s8.71-1.634 12.034-4.902l370.151-363.866v740.062c0 9.24 7.621 16.731 17.02 16.731s17.02-7.492 17.02-16.731V168.13l362.06 355.911c3.324 3.268 7.679 4.902 12.034 4.902s8.71-1.634 12.034-4.902a16.523 16.523 0 0 0 0.001-23.659z"
                        /></svg
                    >
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <svg
                        on:click={() => exchangePosition(i, i + 1)}
                        class:hidden={i + 1 == currentConfig.list.length}
                        class="hover:cursor-pointer fill-gray-500 hover:fill-gray-700 dark:fill-gray-400 dark:hover:fill-white mx-auto mb-1"
                        viewBox="0 0 1024 1024"
                        version="1.1"
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
                        ><path
                            d="M104.8060000000001 523.6179999999999l395.1599999999999 388.45000000000016c6.6489999999999965 6.536000000000002 17.419 6.536000000000004 24.068 7.105427357601002e-15l395.1600000000002-388.4499999999999a16.524 16.524 0 0 0 1.0658141036401503e-14-23.659c-3.3239999999999985-3.2680000000000007-7.6789999999999985-4.902000000000002-12.033999999999999-4.902000000000004s-8.71 1.6339999999999972-12.034000000000002 4.901999999999996l-370.1510000000001 363.8659999999999 2.2737367544323206e-13-740.0619999999999c2.6645352591003757e-15-9.24-7.6209999999999924-16.731-17.019999999999992-16.73100000000001s-17.020000000000003 7.491999999999997-17.020000000000007 16.730999999999995L490.93499999999995 855.8699999999999l-362.05999999999983-355.9110000000001c-3.3239999999999985-3.2680000000000007-7.6789999999999985-4.902000000000002-12.033999999999999-4.902000000000004s-8.71 1.6339999999999972-12.034000000000002 4.901999999999996a16.523 16.523 0 0 0-0.0010000000000083276 23.659z"
                        /></svg
                    >

                    {#if outputDataIndex == i}
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <svg
                            class="hover:cursor-pointer mx-auto mb-1 fill-gray-700 dark:fill-white"
                            viewBox="0 0 1024 1024"
                            version="1.1"
                            xmlns="http://www.w3.org/2000/svg"
                            width="22"
                            height="22"
                            on:click={() => onlyShowOutput(i)}
                            ><path
                                d="M942.2 486.2C847.4 286.5 704.1 186 512 186c-192.2 0-335.4 100.5-430.2 300.3-7.7 16.2-7.7 35.2 0 51.5C176.6 737.5 319.9 838 512 838c192.2 0 335.4-100.5 430.2-300.3 7.7-16.2 7.7-35 0-51.5zM512 766c-161.3 0-279.4-81.8-362.7-254C232.6 339.8 350.7 258 512 258c161.3 0 279.4 81.8 362.7 254C791.5 684.2 673.4 766 512 766z"
                            /><path
                                d="M508 336c-97.2 0-176 78.8-176 176s78.8 176 176 176 176-78.8 176-176-78.8-176-176-176z m0 288c-61.9 0-112-50.1-112-112s50.1-112 112-112 112 50.1 112 112-50.1 112-112 112z"
                            /></svg
                        >
                    {:else}
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <svg
                            class="hover:cursor-pointer mx-auto mb-1 fill-gray-500 hover:fill-gray-700 dark:fill-gray-400 dark:hover:fill-white"
                            viewBox="0 0 1024 1024"
                            version="1.1"
                            xmlns="http://www.w3.org/2000/svg"
                            width="22"
                            height="22"
                            on:click={() => onlyShowOutput(i)}
                            ><path
                                d="M93.866667 322.773333a8.533333 8.533333 0 0 1 6.613333 3.114667c5.589333 6.848 10.261333 12.373333 14.058667 16.64 97.664 109.056 239.552 177.706667 397.482666 177.706667 162.752 0 308.48-72.917333 406.314667-187.84 1.493333-1.792 3.242667-3.882667 5.184-6.272a8.533333 8.533333 0 0 1 15.146667 5.376v9.813333l0.021333 8.32v51.754667a8.533333 8.533333 0 0 1-2.517333 6.037333 599.893333 599.893333 0 0 1-99.584 81.002667l82.474666 98.261333a8.533333 8.533333 0 0 1-1.066666 12.010667l-35.946667 30.165333a8.533333 8.533333 0 0 1-12.010667-1.045333l-89.813333-107.050667a593.045333 593.045333 0 0 1-145.450667 50.837333l44.074667 121.024a8.533333 8.533333 0 0 1-5.098667 10.944l-44.096 16.042667a8.533333 8.533333 0 0 1-10.944-5.098667l-48.448-133.098666a604.586667 604.586667 0 0 1-130.88-1.557334L390.4 714.517333a8.533333 8.533333 0 0 1-10.944 5.12l-44.096-16.064a8.533333 8.533333 0 0 1-5.12-10.944l45.184-124.096a593.066667 593.066667 0 0 1-131.584-47.744l-89.813333 107.029334a8.533333 8.533333 0 0 1-12.032 1.066666L106.026667 598.677333a8.533333 8.533333 0 0 1-1.066667-12.010666l82.474667-98.261334a599.872 599.872 0 0 1-80.981334-62.976c-4.352-4.032-10.56-10.026667-18.602666-18.005333A8.533333 8.533333 0 0 1 85.333333 401.386667v-70.101334c0-4.693333 3.84-8.533333 8.533334-8.533333z"
                            /></svg
                        >
                    {/if}
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <svg
                        on:click={() => deleteConditionAndTemplate(i)}
                        class="hover:cursor-pointer hover:opacity-100 opacity-60 mx-auto mb-1"
                        viewBox="0 0 1024 1024"
                        version="1.1"
                        xmlns="http://www.w3.org/2000/svg"
                        width="22"
                        height="22"
                        ><path
                            d="M768 384c-19.2 0-32 12.8-32 32l0 377.6c0 25.6-19.2 38.4-38.4 38.4L326.4 832c-25.6 0-38.4-19.2-38.4-38.4L288 416C288 396.8 275.2 384 256 384S224 396.8 224 416l0 377.6c0 57.6 44.8 102.4 102.4 102.4l364.8 0c57.6 0 102.4-44.8 102.4-102.4L793.6 416C800 396.8 787.2 384 768 384z"
                            fill="#d81e06"
                        /><path
                            d="M460.8 736l0-320C460.8 396.8 448 384 435.2 384S396.8 396.8 396.8 416l0 320c0 19.2 12.8 32 32 32S460.8 755.2 460.8 736z"
                            fill="#d81e06"
                        /><path
                            d="M627.2 736l0-320C627.2 396.8 608 384 588.8 384S563.2 396.8 563.2 416l0 320C563.2 755.2 576 768 588.8 768S627.2 755.2 627.2 736z"
                            fill="#d81e06"
                        /><path
                            d="M832 256l-160 0L672 211.2C672 166.4 633.6 128 588.8 128L435.2 128C390.4 128 352 166.4 352 211.2L352 256 192 256C172.8 256 160 268.8 160 288S172.8 320 192 320l640 0c19.2 0 32-12.8 32-32S851.2 256 832 256zM416 211.2C416 198.4 422.4 192 435.2 192l153.6 0c12.8 0 19.2 6.4 19.2 19.2L608 256l-192 0L416 211.2z"
                            fill="#d81e06"
                        /></svg
                    >
                </div>
            </div>
            <div class="h-2" />
        {/each}

        <div class="flex justify-center">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <svg
                on:click={addConditionAndTemplate}
                class="hover:cursor-pointer"
                viewBox="0 0 1024 1024"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                width="32"
                height="32"
                ><path
                    d="M480 64A416.64 416.64 0 0 0 64 480 416.64 416.64 0 0 0 480 896 416.64 416.64 0 0 0 896 480 416.64 416.64 0 0 0 480 64z m0 64C674.752 128 832 285.248 832 480S674.752 832 480 832A351.552 351.552 0 0 1 128 480C128 285.248 285.248 128 480 128zM448 320v128H320v64h128v128h64V512h128V448H512V320z"
                    fill="#8a8a8a"
                /></svg
            >
        </div>
    </div>

    <!-- 配置文件列表 -->
    <div class="w-full" class:hidden={!isList}>
        <div
            class="bg-white flex items-center w-full dark:text-gray-400 dark:bg-gray-900"
        >
            <div>
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <svg
                    on:click={() => {
                        isList = false;
                    }}
                    class="hover:cursor-pointer"
                    viewBox="0 0 1024 1024"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    width="16"
                    height="16"
                    ><path
                        d="M670.977781 808.954249c-5.300726 0-10.596336-2.045589-14.603603-6.126534L368.69006 509.86743c-7.818059-7.961322-7.818059-20.717857 0-28.67918l287.684118-292.960285c7.92039-8.065699 20.877493-8.182356 28.942169-0.26299 8.065699 7.919367 8.182356 20.877493 0.264013 28.942169L411.976936 495.526817l273.603425 278.620695c7.918343 8.064676 7.801686 21.022803-0.264013 28.942169C681.331593 807.002804 676.153664 808.954249 670.977781 808.954249z"
                        fill="#8a8a8a"
                    /></svg
                >
            </div>
            <span class="text-xs p-2">Select Config</span>
        </div>
        <div class="bg-white w-full dark:bg-gray-900">
            {#each allConfig.files as r}
                <div class="flex w-full items-center">
                    {#if r.id == allConfig.current}
                        <button
                            on:click={() => {
                                selectConfig(r.id);
                            }}
                            aria-current="true"
                            type="button"
                            class="block w-full cursor-pointer rounded-lg bg-indigo-500 text-white p-4 text-left text-primary-600"
                        >
                            {r.name}
                        </button>
                    {:else}
                        <button
                            on:click={() => {
                                selectConfig(r.id);
                            }}
                            type="button"
                            class="block w-full cursor-pointer rounded-lg p-4 text-left transition duration-500 hover:bg-neutral-100 hover:text-neutral-500 focus:bg-neutral-100 focus:text-neutral-500 focus:ring-0 dark:hover:bg-neutral-600 dark:hover:text-neutral-200 dark:focus:bg-neutral-600 dark:focus:text-neutral-200 dark:bg-gray-700 dark:text-gray-400"
                        >
                            {r.name}
                        </button>
                    {/if}

                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <svg
                        on:click={() => {
                            deleteConfig(r.id);
                        }}
                        class="hover:cursor-pointer mx-2"
                        viewBox="0 0 1024 1024"
                        version="1.1"
                        xmlns="http://www.w3.org/2000/svg"
                        width="22"
                        height="22"
                        ><path
                            d="M768 384c-19.2 0-32 12.8-32 32l0 377.6c0 25.6-19.2 38.4-38.4 38.4L326.4 832c-25.6 0-38.4-19.2-38.4-38.4L288 416C288 396.8 275.2 384 256 384S224 396.8 224 416l0 377.6c0 57.6 44.8 102.4 102.4 102.4l364.8 0c57.6 0 102.4-44.8 102.4-102.4L793.6 416C800 396.8 787.2 384 768 384z"
                            fill="#d81e06"
                        /><path
                            d="M460.8 736l0-320C460.8 396.8 448 384 435.2 384S396.8 396.8 396.8 416l0 320c0 19.2 12.8 32 32 32S460.8 755.2 460.8 736z"
                            fill="#d81e06"
                        /><path
                            d="M627.2 736l0-320C627.2 396.8 608 384 588.8 384S563.2 396.8 563.2 416l0 320C563.2 755.2 576 768 588.8 768S627.2 755.2 627.2 736z"
                            fill="#d81e06"
                        /><path
                            d="M832 256l-160 0L672 211.2C672 166.4 633.6 128 588.8 128L435.2 128C390.4 128 352 166.4 352 211.2L352 256 192 256C172.8 256 160 268.8 160 288S172.8 320 192 320l640 0c19.2 0 32-12.8 32-32S851.2 256 832 256zM416 211.2C416 198.4 422.4 192 435.2 192l153.6 0c12.8 0 19.2 6.4 19.2 19.2L608 256l-192 0L416 211.2z"
                            fill="#d81e06"
                        /></svg
                    >
                </div>
            {/each}
        </div>
    </div>
</div>

<style>
</style>
