package annotate

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Order struct {
	Total   int64
	Item    string
	Id      string
	At      time.Time
	Shipped time.Time
}

type Orders struct {
	Orders  []Order
	byPrice map[int64][]Order
	Log     *zap.Logger
}

func NewOrders() *Orders {
	s := Orders{
		byPrice: map[int64][]Order{},
		Log:     zap.NewNop(),
	}

	return &s
}

func (s *Orders) AddOrder(o Order) {
	s.Log.Debug("Add map entry for", zap.Int64("price", o.Total), zap.Time("at", o.At))
	s.byPrice[o.Total] = append(s.byPrice[o.Total], o)
	s.Orders = append(s.Orders, o)
}

func (s *Orders) ExtractAt(price int64, start time.Time) (Order, error) {
	s.Log.Debug("Looking for", zap.Int64("price", price),
		zap.Time("time", start))
	found := Order{}
	orders, exists := s.byPrice[price]
	if !exists {
		s.Log.Debug("No map entry for", zap.Int64("price", price))
		return found, fmt.Errorf("No order for price %d", price)
	}
	for i, order := range orders {
		if order.At.After(start) || order.At.Equal(start) {
			found = order
			s.byPrice[price] = append(s.byPrice[price][:i], s.byPrice[price][i+1:]...)
			return found, nil
		}
	}

	return found, fmt.Errorf("No order at/after %s", start)
}

func (s *Orders) ParseHistoryFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// skip header bytes
	_, err = file.Seek(3, os.SEEK_CUR)
	if err != nil {
		return err
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	line := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line += 1
		if line == 1 {
			continue
		}

		if len(record) > 23 {
			at, err := timeParse(record[2])
			if err != nil {
				s.Log.Error("Can't parse",
					zap.String("order time", record[2]))
			}
			shipped, err := timeParse(record[18])
			if err != nil {
				s.Log.Error("Can't parse",
					zap.String("ship time", record[18]))
			}

			amt, err := IntFromAmount(record[9])
			if err != nil {
				err = fmt.Errorf("Line %d: %s", line, err)
				return err
			}

			entry := Order{
				Item:    record[23],
				At:      at,
				Shipped: shipped,
				Id:      record[1],
				Total:   amt,
			}

			s.AddOrder(entry)
		} else {
			return fmt.Errorf("Not enough fields: %v", record)
		}
	}

	return nil
}

func asAmount(v int64) string {
	return fmt.Sprintf("%.2f", float64(v)/100)
}

func IntFromAmount(s string) (int64, error) {
	found := int64(0)

	s = strings.ReplaceAll(s, "$", "")

	fl, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", ""), 64)
	if err != nil {
		return found, err
	}
	return int64(fl * 100), nil
}

func timeParse(s string) (time.Time, error) {
	spLoc := strings.LastIndex(s, " ")
	if spLoc != -1 {
		s = s[spLoc+1:]
	}

	dt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}

	return dt, nil
}
