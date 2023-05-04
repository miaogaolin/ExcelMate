<script>
    import { outputData } from "./lib/store.js";
    let data = {};
    let value = "";

    function updateData(data) {
        value = "";
        data.list.forEach((v) => {
            if (data.config_index >= 0 && data.config_index != v.config_index) {
                return;
            }
            value += v.text + "\n";
        });
    }
    outputData.subscribe((v) => {
        data = v;
        updateData(data);
    });
</script>

<div class="result-container h-1/2 w-full dark:bg-gray-900">
    <textarea
        class="h-full placeholder:italic placeholder:text-slate-400 block bg-white w-full border border-slate-300 py-2 pl-2.5 pr-3 outline-none sm:text-sm dark:bg-gray-900 dark:text-gray-400 dark:border-gray-700 tracking-wider bg-gray-200"
        placeholder="Result">{value}</textarea
    >
</div>
