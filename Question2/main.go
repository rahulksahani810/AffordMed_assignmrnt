package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// Product represents a product.
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Company  string  `json:"company"`
	Category string  `json:"category"`
	Rating   float64 `json:"rating"`
	Discount float64 `json:"discount"`
}

// TopProductsResponse represents the response for top products API.
type TopProductsResponse struct {
	Products []Product `json:"products"`
}

// ProductDetailsResponse represents the response for product details API.
type ProductDetailsResponse struct {
	Product Product `json:"product"`
}

func main() {
	http.HandleFunc("/categories/", handleCategories)
    http.HandleFunc("/categories/{categoryname}/products/", handleProductDetails)


	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCategories(w http.ResponseWriter, r *http.Request) {
	categoryName := getCategoryName(r.URL.Path)
	queryParams := r.URL.Query()

	n, err := strconv.Atoi(queryParams.Get("n"))
	if err != nil || n <= 0 {
		n = 10 // Default value
	}

	minPrice, _ := strconv.ParseFloat(queryParams.Get("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(queryParams.Get("maxPrice"), 64)

	// Fetch products from each company
	companies := []string{"AMZ", "FLP", "SNP", "MYN", "AZO"}
	var allProducts []Product
	for _, company := range companies {
		products, err := fetchProducts(company, categoryName, n, minPrice, maxPrice)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching products from %s: %s", company, err.Error()), http.StatusInternalServerError)
			return
		}
		allProducts = append(allProducts, products...)
	}

	// Sort products based on query parameters
	sortProducts(allProducts, queryParams)

	// Return top N products
	response := TopProductsResponse{
		Products: allProducts[:n],
	}
	jsonResponse(w, response)
}

func handleProductDetails(w http.ResponseWriter, r *http.Request) {
	categoryName, productID := getCategoryName(r.URL.Path), getProductID(r.URL.Path)

	// Fetch product details from Test Server
	product, err := fetchProductDetails(categoryName, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching product details: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := ProductDetailsResponse{
		Product: product,
	}
	jsonResponse(w, response)
}

func getCategoryName(path string) string {
	segments := splitPath(path)
	return segments[1]
}

func getProductID(path string) string {
	segments := splitPath(path)
	return segments[3]
}

func splitPath(path string) []string {
	return removeEmptyStrings(strings.Split(path, "/"))
}

func removeEmptyStrings(strings []string) []string {
	var result []string
	for _, s := range strings {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func fetchProducts(company, category string, n int, minPrice, maxPrice float64) ([]Product, error) {
	// Mocking fetching products from the test server
	// In a real scenario, call the actual API of the respective e-commerce company
	var products []Product

	// Simulate fetching products from the test server
	testServerURL := fmt.Sprintf("http://test-server.com/products?company=%s&category=%s&minPrice=%f&maxPrice=%f&n=%d", company, category, minPrice, maxPrice, n)
	resp, err := http.Get(testServerURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching products from test server: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching products from test server: status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, fmt.Errorf("error decoding products from test server response: %s", err.Error())
	}

	return products, nil
}

func fetchProductDetails(category, productID string) (Product, error) {
	// Mocking fetching product details from the test server
	// In a real scenario, call the actual API of the respective e-commerce company

	// Simulate fetching product details from the test server
	testServerURL := fmt.Sprintf("http://test-server.com/products/%s/details", productID)
	resp, err := http.Get(testServerURL)
	if err != nil {
		return Product{}, fmt.Errorf("error fetching product details from test server: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Product{}, fmt.Errorf("error fetching product details from test server: status code %d", resp.StatusCode)
	}

	var product Product
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return Product{}, fmt.Errorf("error decoding product details from test server response: %s", err.Error())
	}

	return product, nil
}

func sortProducts(products []Product, queryParams url.Values) {
	sortBy := queryParams.Get("sortBy")
	sortOrder := queryParams.Get("sortOrder")

	switch sortBy {
	case "rating":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Rating < products[j].Rating
			} else {
				return products[i].Rating > products[j].Rating
			}
		})
	case "price":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Price < products[j].Price
			} else {
				return products[i].Price > products[j].Price
			}
		})
	case "company":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Company < products[j].Company
			} else {
				return products[i].Company > products[j].Company
			}
		})
	case "discount":
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Discount < products[j].Discount
			} else {
				return products[i].Discount > products[j].Discount
			}
		})
	default:
		// Default sorting by name
		sort.Slice(products, func(i, j int) bool {
			if sortOrder == "asc" {
				return products[i].Name < products[j].Name
			} else {
				return products[i].Name > products[j].Name
			}
		})
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
