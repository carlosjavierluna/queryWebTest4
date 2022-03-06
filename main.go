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
	tipo  string
	par   int
	impar int
}

func main() {

	// se definen un canal bidireccional
	canal := make(chan int)
	// Se declara un arreglo de estructuras
	var numeros []numx

	//var i int // este es el valor que se va a guardar en la estructura.

	fmt.Println("Numeros antes del lazo for: ", numeros) // muestra como se creo el arreglo de numeros.

	// Se puede llamar a la funcion llenar_arreglos
	// antes ya que la funcion espera a que
	// el canal tenga datos.
	go llenar_arreglo(&numeros, canal)

	// USamos la go rutina para que los numeros se guarden
	// en el canal en cualquier orden.
	for i := 0; i < 20; i++ {
		go llenar_canal(canal, i)
	}

	time.Sleep(5 * time.Millisecond)          // Se pone este delay para esperar que terminen todas las go rutinas.
	fmt.Println("Despues del for: ", numeros) // se imprime el slice de numeros luego de llamar a las go rutinas.
	texto := fmt.Sprintf("Largo del arreglo: %d Capacidad del arreglo: %d", len(numeros), cap(numeros))
	fmt.Println(texto)
}

//--LLENAR CANAL -------------------------------------------
//  Esta funcion llena el numero que le llega como parametro
//  y lo mete dentro del canal.
func llenar_canal(ch chan<- int, num int) {
	ch <- num
}

//--LLENAR ARREGLO -----------------------------------------
// Esta es la funcion que recibe todo el slice (puntero)
// y le aniade una nuevo elemento considerando si el parametro
// que obtiene del canal  valor de (i)
// es par o impar.
func llenar_arreglo(x *[]numx, ch <-chan int) { // la declaracion es un puntero al un slice de numx
	//fmt.Println("Valor de i:",i)
	for {
		i := <-ch
		// Se determina si es par o impar
		if i%2 == 0 {
			*x = append(*x, numx{"par", 0, i})
		} else {
			*x = append(*x, numx{"Impar", i, 0})
		}

	}
}
