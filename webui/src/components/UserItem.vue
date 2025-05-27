<script>
export default {
    props : ['username']
}
</script>

<template>
    <div class="userItem">
        <!-- Sarebbero le card degli utenti -->
        <button type="button" class="usernameButton" > {{username}} <br> </button>
    </div>
</template>

<style>
.usernameButton{
    font-size: 2rem;
    text-align: center;
    background-color: rgb(106, 160, 177);
    width: 100%;
	font-family: "SF Pro";
	color: #f0f2f3; 
    box-shadow: 7px 5px 5px rgba(4, 103, 133, 1); /* x y sfocatura colore(trasparenza)*/
    border: 1px rgb(4,103,133);
}
.usernameButton:hover{
    background-color: rgb(10, 146, 187);
}
.usernameButton:active{
    box-shadow: 
	inset 0px -12px 12px -6px rgba(31, 77, 97, 0.5), /* x,y,sfocatura,restringimento -> rgba(1 indica l'opacità*/
	inset 0px 12px 12px -6px rgba(11, 51, 68, 0.5),
	inset 12px 0px 12px -6px rgba(11, 51, 68, 0.5),
	inset -12px 0px 12px -6px rgba(11, 51, 68, 0.5);
}
.userItem{
    align-items: center;
	width: 60%;
    margin-bottom: 10px;
}
</style>