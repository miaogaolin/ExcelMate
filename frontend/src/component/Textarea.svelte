<script>
    export let value;
    export let rows;
    export let placeholder;

    let clazz;
    export { clazz as class };

    let textarea;

    function autoResize(el) {
        if (!el) {
            return;
        }
        let target = el;
        if (el.target) {
            target = el.target;
        }
        target.style.height = "auto";

        let height = textarea.scrollHeight;
        if (height > 0) {
            target.style.height = `${height}px`;
        }
    }

    $: autoResize(textarea);
</script>

<textarea
    {rows}
    {value}
    {placeholder}
    class={`${clazz || ""}`}
    bind:this={textarea}
    on:focus={autoResize}
    on:input={autoResize}
    on:keyup
    on:change
/>
