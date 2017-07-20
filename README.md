# dirty code of geofencing experiment
 __  _
 \ `/ |
  \__`!
  / ,' `-.__________________
 '-'\_____                LI`-.
    <____()-=O=O=O=O=O=[]====--)
      `.___ ,-----,_______...-'
           /    .'
          /   .'
         /  .'         upix
         `-'

### Start the web server:

   revel run dev

## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites
