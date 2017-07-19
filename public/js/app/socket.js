var socketManager = function(){
    const JOIN = "join",
        LEAVE = "leave",
        POSITION = "position",
        SETHOOK = "sethook",
        DELHOOK = "delhook",
        RESULT = "result"
        STOPPED = "stopped";

    return {
        socket:null,
        prefixPubSub:"fence",
        getName:function(name){
            if(this.prefixPubSub){
                return this.prefixPubSub+"."+name;
            }else{
                return name;
            }
        },
        init:function(){
            var self=this;
            this.initSocket();
            this.socket.onmessage = function(event){
                var dataStr = event.data;
                var data = JSON.parse(dataStr);
                switch(data.type){
                    case JOIN:;break;
                    case LEAVE:;break;
                    case POSITION:;
                        if(self.onPosition!=null){
                            self.onPosition(data);
                        }
                    break;
                    case SETHOOK:;break;
                    case DELHOOK:;break;
                    case STOPPED:console.log("Stopped ", data);break;
                    case RESULT:
                        if(self.onResult!=null){
                            self.onResult(data);
                        }
                    ;break;
                    default:console.log(event, data);break;
                }
            }
            this.socket.onopen = this.socketConnected;
        },
        initSocket:function(){
            if(!this.socket){
                this.socket = new WebSocket('ws://'+window.location.host+'/websocket/geofence')
            }
        },
        send:function(data){
            this.socket.send(JSON.stringify(data));
        },
        sethook:function(name, group, geojson){
            var self=this;
            var data ={
                geojson: geojson,
                group: group,
                name: self.getName(name),
                cmd: SETHOOK
            }
            this.send(data);
        },
        delhook:function(name){
            var data =  { cmd:DELHOOK, name:name};
            this.send(data);
        },
        updateLocation:function(name, group, lat, long){
            var data =  { cmd:POSITION, name:name, group:group, lat:lat, long:long};
            this.send(data);
        },
        socketConnected:function(fn){
            // var data = {
            //     cmd:"position",
            //     data:"oke bos"
            // };
            // this.send(data);
            if(typeof fn === "function"){
                fn();
            }
        },
        onPosition:null,
        onResult:null,
    }
}
window.socketManager = socketManager;