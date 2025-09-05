<script>
    export default{
        props: ['messageId']
    }
</script>

<template>
    <div class="message-options-overlay">
        <div class="options-overlay-content">
            <div class="emoji-picker">
                <!-- Griglia di emoji da poter inserire-->
                <button @click="$emit('chosen', '👍')"> 👍 </button>
                <button @click="$emit('chosen', '😂')"> 😂 </button>
                <button @click="$emit('chosen', '😢')"> 😢 </button>
                <button @click="$emit('chosen', '😡')"> 😡 </button>
                <button @click="$emit('chosen', '👎')"> 👎 </button>
                 
            </div>
            <div class="buttons">
                <button @click="$emit('close')"> Annulla </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.emoji-picker button{
    padding: 5px 10px;
    background-color: rgb(165, 222, 240);
	border: none;
	border-radius: 8px;
	cursor: pointer;
}
.emoji-picker{
    display: flex;
    align-items: baseline;
    gap: 1.8rem;
}
.buttons button:last-child{
    background-color: rgb(226, 22, 22);
	color: white;
    border-radius: 8px;
    border: none;
    padding: 10px 20px;
    font-family: "SF Pro";
}
.buttons{
    margin-top: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}
.options-overlay-content{
    background: white;
    padding: 2rem;
    border-radius: 12px;
    width: 450px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}
.message-options-overlay{
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.4);
    display: flex;
    align-items: center;
    justify-content: center;
}</style>