<script>
    export default{
        props: ['conversationId', 'users']
    }
</script>
<template>
    <div class="overlay">
        <div class="overlay-content">
            <div class="users-list">
                <ul>
                    <li v-for="user in users" :key="user.userId">
                        <button @click="$emit('addToGroup', user.id)">
                            {{ user.name }}
                        </button>
                    </li>
                </ul>
            </div>
            <button class="cancel-button" @click="$emit('close')">Annulla</button>
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
.users-list button{
    padding: 10px 20px;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-family: "SF Pro";
    width: 100%;
}
.users-list li{
    margin-bottom: 0.5rem;
}
.users-list ul{
    list-style: none;
    padding: 0;
    margin: 0;
}
.users-list{
    flex: 1;
    overflow-y: auto;
    margin-bottom: 1.5rem;
}
.overlay-content{
    background: white;
    padding: 2rem;
    border-radius: 12px;
    width: 450px;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}
.overlay{
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