<script>
    export default{
        props: ['messageId']
    }
</script>

<template>
    <div class="message-options-overlay">
        <div class="options-overlay-content">
            <div class="buttons">
                <button @click="$emit('forward', messageId)"> Inoltra messaggio </button>
                <button @click="$emit('delete',messageId)"> Elimina messaggio </button>
                <button @click="$emit('openEmojiPicker',messageId)"> Commenta messaggio </button>
                <button @click="$emit('openComments',messageId)"> Elimina commento </button>
                <button @click="$emit('close')"> Annulla </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.buttons button:last-child{
    background-color: rgb(226, 22, 22);
	color: white;
}
.buttons button{
    padding: 10px 20px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
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
}
</style>