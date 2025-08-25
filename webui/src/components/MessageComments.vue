<script>
    export default{
        props: ['userId','messageId', 'comments']
    }
</script>

<template>
    <div class="comments-overlay">
        <div class="comments-content">
            <div class="comments-list">
                <ul>
                    <li v-for="comm in comments" :key="comm.commentId">
                        <button @click="$emit('uncomment', comm.commentId)" 
                        :class="{'my-comment': String(comm.authorId) === userId}">
                            {{ comm.emoji }}
                        </button>
                    </li>
                </ul>
            </div>
            <button class="cancel-button" @click="$emit('close')"> Annulla </button>
        </div>
    </div>
</template>

<style scoped>
.cancel-button{
    background-color: rgb(226, 22, 22);
	color: white;
    border-radius: 8px;
    border: none;
    padding: 10px 20px;
    font-family: "SF Pro";
}
.comments-list button.my-comment { 
  border: 1px solid rgb(10, 146, 187);
}
.comments-list button{
    padding: 10px 20px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-family: "SF Pro";
    width: 100%;
}
.comments-list li{
    margin-bottom: 0.5rem;
}
.comments-list ul{
    list-style: none;
    padding: 0;
    margin: 0;
}
.comments-list{
    flex: 1;
    overflow-y: auto;
    margin-bottom: 1.5rem;
}
.comments-content{
    background: white;
    padding: 2rem;
    border-radius: 12px;
    width: 450px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}
.comments-overlay{
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