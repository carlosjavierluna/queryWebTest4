package main

import (
	"fmt"
	"time"
)

// se declara una estructura para almacenar los numeros.
type numx struct {
	tipo  string
	par   int
	impar int
}

func main() {
	var numeros []numx // se declara un slice de la estructura declarada global

	var i int // este es el valor que se va a guardar en la estructura.

	fmt.Println("Numeros antes del lazo for: ", numeros) // muestra como se creo el arreglo de numeros.

	for i = 0; i < 20; i++ { // Se llama a gorutinas para estudiar el comportamiento de las mismas.
		go llenar_numero(&numeros, i)
	}
	time.Sleep(20 * time.Millisecond)         // Se pone este delay para esperar que terminen todas las go rutinas.
	fmt.Println("Despues del for: ", numeros) // se imprime el slice de numeros luego de llamar a las go rutinas.
}

// Esta es la funcion que recibe todo el slice (puntero)
// y le aniade una nuevo elemento considerando si el parametro (i)
// es par o impar.
func llenar_numero(x *[]numx, i int) { // la declaracion es un puntero al un slice de numx
	fmt.Println("Valor de i:", i)
	// Se determina si es par o impar
	if i%2 == 0 {
		*x = append(*x, numx{"par", 0, i})
	} else {
		*x = append(*x, numx{"Impar", i, 0})
	}
}
