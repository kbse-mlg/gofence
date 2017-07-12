var dashboardModule = function($, L){
  var markers = {},s
  layers = [];
  return{
      //attribute
      center:[3.1290962786081646, 101.67658295822144],
      zoom:17,
      id: "gmapid",
      map: null,
      urlTileLayer:"https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw",
      setCenter: function(lat, long){
          this.center = [lat, long];
      },
      init:function(){
          this.map = L.map(this.id, {
              center: this.center,
              zoom: this.zoom
          });

          L.tileLayer(this.urlTileLayer, {
              attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://mapbox.com">Mapbox</a>',
              maxZoom: 18,
              id: 'mapbox.streets',
              accessToken: 'pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw'
          }).addTo(this.map);
      },
      loadArea:function(){
          $.get( "/api/v1/area", function( data ) {
              if(data && data.Areas){
                  data.Areas.forEach(function(area) {
                      var gsonLayer = L.geoJSON(JSON.parse(area.geodata), {
                          style: function (feature) {
                              return {color: feature.properties.color};
                          },
                          onEachFeature: onEachFeature
                      }).bindPopup(function (layer) {
                          return layer.feature.properties.description;
                      });
                      this.layers.push(gsonLayer);
                      gsonLayer.addTo(this.map);
                  }, this);
              }
          });
      },
      loadObject:function(){
          $.get( "/api/v1/object", function( data ) {
              if(data && data.Objects){
                  data.Objects.forEach(function(obj) {
                      markers[obj.Name] = L.marker([obj.lat, obj.long]);
                      markers[obj.Name].bindPopup(obj.name).addTo(this.map)
                  });
              }
          });
      }
  };
}( jQuery , L);
