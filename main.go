package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

//TODO: Crear un canal de escucha al que le vamos
//      a enviar las estructuras y que cuando las
//      reciba las procese.

// Esta estructura es para almacenar las coordenadas
// de la ciudad
type coordCity struct {
	idxSlice int
	index    int
	state    string
	latCity  float64
	lonCity  float64
}

// Estructura para almacenar los datos del clima de la ciudad
type climaCity struct {
	idxSlice int
	clima    string
}

// Esta estructura es para guardar la informacion que se recupera de internet
// y de los nombres que van como parametros.

type dataWeb struct {
	index   int
	city    string
	state   string
	country string
	// location string // Se elimina
	latCity float64
	lonCity float64
	weather string
	// Se agregan estos dos campos para saber cuando se crea
	// y cuando se termina de llenar los datos.
	tcreate time.Time // cuando se crea el dato
	tdone   time.Time // cuando finaliza la creacion

}

// Estructura para almacenar los nombres de las ciudades
// Que tienen mal formato en la linea de comandos.
type dataBad struct {
	index int
	city  string
}

// Estructura para recuperar en formato JSON
// la geo localizacion de la ciudad.

type Ubicacion struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

// Estructura para recupear en formato JSON
// el clima de una ciudad.

// la INFORMACION DEL CLIMA EN UNA CIUDAD
type coord struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}
type weather struct {
	Id          int
	Main        string
	Description string
	Icon        string
}
type Mein struct {
	Temp       float64 `json:"temp"`
	Feels_like float64 `json:"feels_like"`
	Temp_min   float64 `json:"temp_min"`
	Temp_max   float64 `json:"temp_max"`
	Pressure   int     `json:"pressure"`
	Humidity   int     `json:"humidity"`
}

type wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

type clouds struct {
	All int `json:"all"`
}

type sys struct {
	Type    int    `json:"type"`
	Id      int16  `json:"id"`
	Country string `json:"country"`
	Sunrise int32  `json:"sunrise"`
	Sunset  int32  `json:"sunset"`
}

type Clima struct {
	Coord      coord   `json:"coord"`
	Weather    weather `json:"weather"`
	Base       string  `json:"base"`
	Main       Mein    `json:"main"`
	Visibility int     `json:"visibility"`
	Wind       wind    `json:"wind"`
	Clouds     clouds  `json:"clouds"`
	Dt         int32   `json:"dt"`
	Sys        sys     `json:"sys"`
	Timezone   int32   `json:"timezone"`
	Id         int32   `json:"id"`
	Name       string  `json:"name"`
	Cod        int     `json:"cod"`
}

//--MAIN-------------------------------------------

func main() {

	// Estas variables contienen los arreglos con parametros buenos y malos
	var dataCity []dataWeb // Ciudades con formato OK
	var badCity []dataBad  // Ciudades con formato no OK

	// 1. Se leen y procesan los parametros de entrada.
	// Se ingresan los args es un slice de strings
	args := os.Args
	//fmt.Println(reflect.TypeOf(args)) // Se obtiene el tipo de datos de una variable

	// Si no hay argumentos finaliza el programa.
	if len(args[1:]) == 0 {
		fmt.Println("Se debe ingresar al menos un  nombre de Ciudad...")
		return
	}

	// La siguiente funcion verifica los nombres de ciudades
	// pasados como parametros de la linea de comandos y los ingresa
	// en los arreglos correspondientes.
	//
	// Al salir de la funcion de verificar ciudades los slices contienen los datos

	verificar_ciudades(&dataCity, &badCity, args)

	// Se muestra un mensaje del numero de argumentos con formato incorrecto
	if len(badCity) > 0 {
		fmt.Println("\n\nAlgunos argumentos (", len(badCity), ") tienen formato incorrecto...")
		fmt.Println(badCity)
	}
	// Se definen los canales donde se pondran los datos.
	canal_ciudades := make(chan dataWeb)
	canal_coordenadas := make(chan coordCity)
	s := make(chan int)

	fmt.Println("\nINICIA EL PROCESAMIENTO WEB\n")
	// Inicia la medicion del tiempo de ejecucion luego
	// de que se selecciona la opcion.
	now := time.Now()
	// Esta funcion envia los nombres de las ciudades al canal.
	go enviar_ciudades(dataCity, canal_ciudades, s)
	// Dentro de la funciion se obtienen las coordenadas (lat y lon) de la ciudad
	// El arreglo_coordenadas contiene el resultado de la funcion.
	arreglo_coordenadas := obtener_coordenadas(canal_ciudades, s)
	// Se envia el slice de las coordenadas para obtener el clima
	// en las coordenadas.
	go enviar_coordenadas(arreglo_coordenadas, canal_coordenadas, s)
	// Obtiene el clima para cada par de coordenadas y devuelve un Slice
	// arreglo_climas donde estan todos los datos del clima formateados.
	arreglo_climas := obtener_clima(canal_coordenadas, s)
	// En el arreglo de ciudades se actualiza las coordenadas (Lat y Lon)
	// y el string del clima.
	for i := 0; i < len(arreglo_climas); i++ {
		// Se puede utilizar el mismo contador por que los dos slices
		// siempre tienen la misma longitud.
		dataCity[arreglo_coordenadas[i].idxSlice].state = arreglo_coordenadas[i].state
		dataCity[arreglo_coordenadas[i].idxSlice].latCity = arreglo_coordenadas[i].latCity
		dataCity[arreglo_coordenadas[i].idxSlice].lonCity = arreglo_coordenadas[i].lonCity
		dataCity[arreglo_coordenadas[i].idxSlice].tdone = time.Now()

		dataCity[arreglo_climas[i].idxSlice].weather = arreglo_climas[i].clima
	}
	// Se imprimen los resultados obtenidos
	impCityes(dataCity)
	fmt.Println("\n\nTiempo transcurrido:", time.Since(now))
	return // Finaliza el programa
}

//--ENVIAR_CIUDADES: Envia las ciudades a un canal para se procesadas una a una.
func enviar_ciudades(ciudades []dataWeb, canal_ciudades chan<- dataWeb, s chan<- int) {
	//fmt.Println("\n\nEmpieza ENVIAR_CIUDADES.... \n\n")
	for i := 0; i < len(ciudades); i++ {
		canal_ciudades <- ciudades[i]
	}
	s <- 0 // Le permite a la rutina de obtencion terminar.
}

//--ENVIAR COORDENADAS: envia las coordenadas de cada ciudad para
//  obtener luego el clima de esas coordenadas
func enviar_coordenadas(arrayCoordenadas []coordCity, canal_coordenadas chan<- coordCity, s chan<- int) {
	for i := 0; i < len(arrayCoordenadas); i++ {
		canal_coordenadas <- arrayCoordenadas[i]
	}
	s <- 0 // Le permite a la rutina de obtencion terminar.
}

//--OBTENER COORDENADAS------------------------------------------------
// Obtiene la latitud y la longitud a partir del nombre de la ciudad y pais
// en el formato Ciudad PA
func obtener_coordenadas(ciudades <-chan dataWeb, s <-chan int) []coordCity {
	var arrayCoordenadas []coordCity
	posicionArreglo := 0

	for {
		select {
		case v := <-ciudades:
			client := &http.Client{Timeout: 30 * time.Second}
			// GEOLOCALIZACION
			urlFind := "https://api.openweathermap.org/geo/1.0/direct?q=" + v.city + "," + v.country + "&limit=1&appid=3582889a6bd6c9bd3d2867e51f420116"
			bytes, err := getWebBytes(client, urlFind)
			// manejo de error
			if err != nil {
				log.Fatal(err)
			}
			// Almacena la ubicacion recuperada de la ciudad
			// en la llamada a la funcion queryCityLocation
			var cityLoc Ubicacion
			//--APLICANDO UNMARSHAL-------------------------
			// bytes recupera [] al inicio y final, no se toman en cuenta.
			objetostring := string(bytes)[1 : len(bytes)-1]
			json.Unmarshal([]byte(objetostring), &cityLoc)

			arrayCoordenadas = append(arrayCoordenadas, coordCity{posicionArreglo, v.index, cityLoc.State, float64(cityLoc.Lat), float64(cityLoc.Lon)})
			posicionArreglo = posicionArreglo + 1 // Se incrementa el indice
		case v := <-s:
			fmt.Println("Se obtuvo Lat y Lon de las ciudades...", v)
			return arrayCoordenadas // Se sale de la funcion cuando se recibe en este canal... OK
		}
	}
} //--FIN--OBTENER COORDENADAS------------------------------------------------

//--OBTENER CLIMA----------------------------------------------------------
// Obtiene el clima de una localidad a partir de la longitud y latitud del
// lugar.
func obtener_clima(coordenadas <-chan coordCity, s <-chan int) []climaCity {
	var arrayClima []climaCity
	for {
		select {
		case v := <-coordenadas:
			//fmt.Println("\n\n\nDesde el canal de CLIMA  : ", v, "\n\n")
			client := &http.Client{Timeout: 30 * time.Second}

			// CLIMA
			var cadena float64
			cadena = v.lonCity
			strLongitud := strconv.FormatFloat(cadena, 'f', 7, 64)
			strLatitud := strconv.FormatFloat(v.latCity, 'f', 7, 64)

			urlFind := "https://api.openweathermap.org/data/2.5/weather?units=metric&lat=" + strLatitud + "&lon=" + strLongitud + "&appid=3582889a6bd6c9bd3d2867e51f420116"
			//fmt.Println("\nURL PARA OBTENER EL CLIMA:\n\n", urlFind, "\n\n")
			bytes, err := getWebBytes(client, urlFind)

			// manejo de error
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println("Dados de Clima: ")
			// fmt.Println(string(bytes))

			var vClima Clima // Esta variable es para almacenar en formato JSON lo recuperado en la consulta web.

			//--APLICANDO UNMARSHAL-------------------------")
			objetostring := string(bytes)
			json.Unmarshal([]byte(objetostring), &vClima)

			// fmt.Println("\n\nDESPUES DE APLICAAR EL SEGUNDO UNMARSHAL")
			// fmt.Println(vClima)

			wproceced := fmt.Sprintf("\n\tTemperatura   : %2.2f Grados", vClima.Main.Temp)
			//wproceced = wproceced + fmt.Sprintf("\n\tSensacion Term: %f", vClima.Main.Feels_like)
			//wproceced = wproceced + fmt.Sprintf("\n\tTemp Mínima   : %f", vClima.Main.Temp_min)
			wproceced = wproceced + fmt.Sprintf("\n\tTemp Máxima   : %2.2f Grados", vClima.Main.Temp_max)
			wproceced = wproceced + fmt.Sprintf("\n\tPresión Atm.  : %dhPa", vClima.Main.Pressure)
			wproceced = wproceced + fmt.Sprintf("\n\tHumedad       : %d %%", vClima.Main.Humidity)
			wproceced = wproceced + fmt.Sprintf("\n\tVisibilidad   : %d mts", vClima.Visibility)

			//fmt.Println("CADENA APLICADA JSON: ", wproceced)

			arrayClima = append(arrayClima, climaCity{v.idxSlice, wproceced})

		case v := <-s:
			fmt.Println("Se obtuvo el clima de las ciudades ...", v)
			return arrayClima // Se sale de la funcion cuando se recibe en este canal... OK
		}
	}
} //--FIN--OBTENER CLIMA----------------------------------------------------------

//--GETWEBBYTES (QUERYWEB) -------------------------------------------
func getWebBytes(client *http.Client, url string) ([]byte, error) {

	// Se verifica si la solicitud se ha construido bien.
	// Si se tiene un error en la conexion o la consulta se termina el programa
	// y se muestra el mensaje de error generado.

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Ejecutar la consulta. El método Do ejecuta la soliciutd.
	res, err := client.Do(req)
	// Si se tiene un error en la conexion o la consulta se termina el programa
	// y se muestra el mensaje de error generado.
	if err != nil {
		log.Fatal(err)
	}

	// La siguiente funcion lee todo el "body" de la pagina web.
	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	return bytes, nil

} //FIN--GETWEBBYTES (QUERYWEB) -------------------------------------------

//--VERIFICAR CIUDADES-----------------------------------------------------
// Verifica que los parametros de ciudades cumplan con el formato Ciudad,PA
// Si hay nombres mal escritos los almacena en otro arreglo.
func verificar_ciudades(ciudades *[]dataWeb, malas *[]dataBad, argumentos []string) {
	// 2. Se verifica que los nombres de las ciudades
	// sean del formato ciudad,PA donde PA son las siglas del pais.

	//crear una variable global que almacenará la expresión regular
	// Se modifica la expresion regular para que acepte tildes el el nombre
	expresionRegular := regexp.MustCompile("^[A-Z][A-Za-záéíóú]+[,]+([A-Z][A-Z])?$")

	for idx, arg := range argumentos[1:] {
		// range va desde 0 hasta el tamanio del arreglo
		if expresionRegular.Match([]byte(arg)) { // Si coincide con la expresion regular es TRUE
			// Se descompone la ciudad en nombre y pais Quito , EC
			posComma := strings.Index(arg, ",")
			nameCity := arg[0:posComma]
			codCountry := arg[posComma+1:]
			// 3. Se guardan los parametros ok en un arreglos de ciudades
			*ciudades = append(*ciudades, dataWeb{idx + 1, nameCity, "vacio", codCountry, 0, 0, "vacio", time.Now(), time.Now()})
		} else {
			// Se guardan las ciudades con formato incorrecto
			// en un arreglo de malas ciudades
			*malas = append(*malas, dataBad{idx + 1, arg})
		}

	}

} //--FIN--VERIFICAR CIUDADES-----------------------------------------------------

//--IMPCITYES---------------------------------------------------------------------------
//
//  Imprime en formato los datos recuperados de la web.
//  para cada una de las ciudades.
func impCityes(cityes []dataWeb) {

	fmt.Println("\n\nCiudades procesadas..\n------------------------------------------------")
	// Se ordena el Slice por el nombre de la ciudad.
	sort.Slice(cityes, func(i, j int) bool {
		return cityes[i].city < cityes[j].city
	})

	for i := 0; i < len(cityes); i++ {
		fmt.Println("No:", i+1)
		fmt.Println("Ciudad      : ", cityes[i].city)
		fmt.Println("Estado      : ", cityes[i].state)
		fmt.Println("Pais        : ", cityes[i].country)

		fmt.Println("Latitud     : ", fmt.Sprintf("%f", cityes[i].latCity))
		fmt.Println("Longitud    : ", fmt.Sprintf("%f", cityes[i].lonCity))
		fmt.Println("Tiempo proc.: ", cityes[i].tdone.Sub(cityes[i].tcreate))
		fmt.Println("Info Clima  : ", cityes[i].weather)

		fmt.Println("------------------------------------------------")
	}
} // END: func impCityes(cityes []dataWeb)
