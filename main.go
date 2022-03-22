package main

import (
	"fmt"
	"time"
)

//TODO: Crear un canal de escucha al que le vamos
//      a enviar las estructuras y que cuando las
//      reciba las procese.

// se declara una estructura para almacenar los numeros.
type numx struct {
	idx   int // para guardar el orden en que se van creando
	tipo  string
	par   int
	impar int
}

func main() {
	// se definen un canal bidireccional
	canal := make(chan numx)
	// Se declara un arreglo de estructuras
	var numeros []numx

	//var i int // este es el valor que se va a guardar en la estructura.

	fmt.Println("Numeros antes del lazo for: ", numeros) // muestra como se creo el arreglo de numeros.

	// USamos la go rutina para que los numeros se guarden
	// en el canal en cualquier orden.
	for i := 0; i < 20; i++ {
		go llenar_canal(canal, i)
	}

	// Se puede llamar a la funcion llenar_arreglos
	// antes ya que la funcion espera a que
	// el canal tenga datos.
	go llenar_arreglo(&numeros, canal)

	time.Sleep(10 * time.Millisecond)    // Se pone este delay para esperar que terminen todas las go rutinas.
	fmt.Println("\nDespues del Go...\n") // se imprime el slice de numeros luego de llamar a las go rutinas.
	for i := range numeros {
		fmt.Println("Arreglo [", i, "] ", numeros[i])
	}
	texto := fmt.Sprintf("\nLargo del arreglo: >>> %d <<< Capacidad del arreglo: >>> %d <<<\n", len(numeros), cap(numeros))
	fmt.Println(texto)
}

//--LLENAR CANAL -------------------------------------------
//  Esta funcion llena el numero que le llega como parametro
//  y lo mete la estructura dentro del canal.
func llenar_canal(ch chan<- numx, i int) {
	var x numx // Se declara la variable
	x.idx = i  // Se llena el indice con el valor de i
	ch <- x    //se envia al canal la variable completa
}

//--LLENAR ARREGLO -----------------------------------------
// Esta es la funcion que recibe todo el slice (puntero)
// y le aniade una nuevo elemento considerando si el parametro
// que obtiene del canal  valor de (i)
// es par o impar.
func llenar_arreglo(x *[]numx, ch <-chan numx) {
	// la declaracion es un puntero al un slice de numx
	//fmt.Println("Valor de i:",i)
	for {
		i := <-ch // Se recupera el dato del canal

		fmt.Println("i dentro de llenar_arreglo():", i)

		// Se determina si es par o impar
		if i.idx%2 == 0 {
			i.tipo = "par"
			i.par = i.idx
		} else {
			i.tipo = "Impar"
			i.impar = i.idx
		}
		*x = append(*x, i)

	}
}
