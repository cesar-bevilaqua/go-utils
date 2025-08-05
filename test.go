package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	interval := 500 * time.Millisecond
	filename := "monitoramento.csv"

	// Criar ou abrir o arquivo CSV
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Cabeçalho do CSV
	writer.Write([]string{
		"Timestamp",
		"Uso_CPU_percent",
		"Uso_RAM_percent",
		"Média_CPU",
		"Média_RAM",
		"Pico_CPU",
		"Pico_RAM",
	})

	var totalCPU, totalRAM float64
	var maxCPU, maxRAM float64
	var count int

	fmt.Println("Iniciando monitoramento (Ctrl+C para parar)...")

	for {
		// Coletar uso da CPU (média geral de todos os núcleos)
		cpuPercent, err := cpu.Percent(interval, false)
		if err != nil {
			fmt.Println("Erro ao ler uso da CPU:", err)
			continue
		}

		// Coletar uso de memória
		memStats, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Erro ao ler uso da memória:", err)
			continue
		}

		cpuUsage := cpuPercent[0]
		ramUsage := memStats.UsedPercent

		// Atualizar soma e picos
		totalCPU += cpuUsage
		totalRAM += ramUsage
		if cpuUsage > maxCPU {
			maxCPU = cpuUsage
		}
		if ramUsage > maxRAM {
			maxRAM = ramUsage
		}
		count++

		// Calcular médias
		avgCPU := totalCPU / float64(count)
		avgRAM := totalRAM / float64(count)

		// Timestamp atual
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		// Escrever no CSV
		writer.Write([]string{
			timestamp,
			fmt.Sprintf("%.2f", cpuUsage),
			fmt.Sprintf("%.2f", ramUsage),
			fmt.Sprintf("%.2f", avgCPU),
			fmt.Sprintf("%.2f", avgRAM),
			fmt.Sprintf("%.2f", maxCPU),
			fmt.Sprintf("%.2f", maxRAM),
		})

		writer.Flush() // Escreve imediatamente no disco

		// Print no terminal (opcional)
		fmt.Printf("[%s] CPU: %.2f%% (média: %.2f, pico: %.2f) | RAM: %.2f%% (média: %.2f, pico: %.2f)\n",
			timestamp, cpuUsage, avgCPU, maxCPU, ramUsage, avgRAM, maxRAM)
	}
}
