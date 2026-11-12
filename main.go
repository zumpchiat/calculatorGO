package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Calculadora Go")
	w.Resize(fyne.NewSize(260, 320))
	w.SetFixedSize(true)

	// Define as cores (R, G, B, Alpha)
	corNormal := color.NRGBA{R: 0, G: 220, B: 100, A: 255} // Verde
	corErro := color.NRGBA{R: 255, G: 70, B: 70, A: 255}   // Vermelho

	// Visor usando canvas.Text para permitir cor personalizada
	visor := canvas.NewText("0", corNormal)
	visor.TextSize = 28
	visor.Alignment = fyne.TextAlignTrailing // Alinha o texto à direita
	visor.TextStyle = fyne.TextStyle{Bold: true}

	// Função auxiliar para atualizar o texto e redesenhar o visor
	atualizarVisor := func(texto string, cor color.Color) {
		visor.Text = texto
		visor.Color = cor
		visor.Refresh() // Obrigatório no canvas.Text para aplicar mudanças na tela
	}

	var numAnterior float64
	var operacao string
	novoNumero := true

	digitar := func(valor string) {
		if novoNumero {
			if valor == "." {
				atualizarVisor("0.", corNormal)
			} else {
				atualizarVisor(valor, corNormal)
			}
			novoNumero = false
		} else {
			if valor == "." && containsDot(visor.Text) {
				return
			}
			if visor.Text == "0" && valor != "." {
				atualizarVisor(valor, corNormal)
			} else {
				atualizarVisor(visor.Text+valor, corNormal)
			}
		}
	}

	definirOperacao := func(op string) {
		val, err := strconv.ParseFloat(visor.Text, 64)
		if err == nil {
			numAnterior = val
			operacao = op
			novoNumero = true
		}
	}

	calcular := func() {
		if operacao == "" {
			return
		}
		numAtual, err := strconv.ParseFloat(visor.Text, 64)
		if err != nil {
			return
		}

		var resultado float64
		switch operacao {
		case "+":
			resultado = numAnterior + numAtual
		case "-":
			resultado = numAnterior - numAtual
		case "*":
			resultado = numAnterior * numAtual
		case "/":
			if numAtual == 0 {
				atualizarVisor("Erro: Div por 0", corErro)
				operacao = ""
				novoNumero = true
				return
			}
			resultado = numAnterior / numAtual
		}

		atualizarVisor(fmt.Sprintf("%g", resultado), corNormal)
		operacao = ""
		novoNumero = true
	}

	limpar := func() {
		atualizarVisor("0", corNormal)
		numAnterior = 0
		operacao = ""
		novoNumero = true
	}

	grade := container.NewGridWithColumns(4,
		widget.NewButton("7", func() { digitar("7") }),
		widget.NewButton("8", func() { digitar("8") }),
		widget.NewButton("9", func() { digitar("9") }),
		widget.NewButton("/", func() { definirOperacao("/") }),

		widget.NewButton("4", func() { digitar("4") }),
		widget.NewButton("5", func() { digitar("5") }),
		widget.NewButton("6", func() { digitar("6") }),
		widget.NewButton("*", func() { definirOperacao("*") }),

		widget.NewButton("1", func() { digitar("1") }),
		widget.NewButton("2", func() { digitar("2") }),
		widget.NewButton("3", func() { digitar("3") }),
		widget.NewButton("-", func() { definirOperacao("-") }),

		widget.NewButton("C", limpar),
		widget.NewButton("0", func() { digitar("0") }),
		widget.NewButton(".", func() { digitar(".") }),
		widget.NewButton("+", func() { definirOperacao("+") }),
	)

	btnIgual := widget.NewButton("=", calcular)
	btnIgual.Importance = widget.HighImportance

	// Adiciona um fundo escuro ao redor do visor para dar destaque à cor da fonte
	fundoVisor := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	caixaVisor := container.NewStack(fundoVisor, container.NewPadded(visor))

	conteudo := container.NewVBox(
		caixaVisor,
		grade,
		btnIgual,
	)

	w.SetContent(conteudo)
	w.ShowAndRun()
}

func containsDot(s string) bool {
	for _, c := range s {
		if c == '.' {
			return true
		}
	}
	return false
}

