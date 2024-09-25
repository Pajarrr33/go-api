package controller

import (
	"net/http"
	"strconv"
	"strings"
	"submission-project-enigma-laundry/entity"
	"submission-project-enigma-laundry/repository"
	"github.com/gin-gonic/gin"
)

type ProductController interface {
	CreateProduct(ctx *gin.Context)
	ListProduct(ctx *gin.Context)
	GetDetailProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
}

type productController struct {
	productRepository repository.ProductRepository
}

type ProductResponse struct {
	Message string `json:"message"`
	Data entity.Product `json:"data"`
}

type ProductResponseSlice struct {
	Message string `json:"message"`
	Data []entity.Product `json:"data"`
}

func NewProductController(repo repository.ProductRepository) ProductController {
	return &productController{productRepository: repo}
}

func (pc *productController) CreateProduct(ctx *gin.Context) {
	var newProduct entity.Product
	err := ctx.ShouldBind(&newProduct)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Invalid Input","details" : err.Error()})
		return
	}

	createdProduct,err := pc.productRepository.CreateProduct(&newProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to create product", "details" : err.Error()})
		return
	}
	// Create The response struct 
	response := ProductResponse{
		Message: "Successfuly Create Product",
		Data: *createdProduct,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (pc *productController) ListProduct(ctx *gin.Context) {
	products := []entity.Product{}
	productName := ctx.Query("productName")
	if strings.TrimSpace(productName) != "" {
		rows, err := pc.productRepository.GetProductByName(productName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get all product by name", "details" : err.Error()})
			return
		}

		defer rows.Close()

		for rows.Next() {
			product := entity.Product{}
			err = rows.Scan(&product.Product_id,&product.Product_name,&product.Price,&product.Unit)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed scanning product data by name", "details" : err.Error()})
				return
			}
			products = append(products, product)
		}

		err = rows.Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error encountred during iteration", "details" : err.Error()})
			return
		}

		response :=ProductResponseSlice{
			Message: "Successfully get data by name from product",
			Data: products,
		}

		ctx.JSON(http.StatusOK, response)
		return
	} else {
		rows, err := pc.productRepository.GetProduct()
	
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get all product data", "details" : err.Error()})
			return
		}

		defer rows.Close()

		for rows.Next() {
			product := entity.Product{}
			err = rows.Scan(&product.Product_id,&product.Product_name,&product.Unit,&product.Price)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed scanning product data", "details" : err.Error()})
				return
			}
			products = append(products, product)
		}

		err = rows.Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error encountred during iteration", "details" : err.Error()})
			return
		}

		response :=ProductResponseSlice{
			Message: "Successfully get all data from product",
			Data: products,
		}

		ctx.JSON(http.StatusOK, response)
		return
	}
}

func (pc *productController) GetDetailProduct(ctx *gin.Context) {
	id,err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message" : "Failed convert id. Make sure id is number", "details" : err.Error()})
		return
	}

	product := entity.Product{}

	detailProduct, err := pc.productRepository.GetDetailProduct(id,&product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Failed to get detail product data", "details" : err.Error()})
		return
	}

	// Create The response struct 
	response := ProductResponse{
		Message: "Successfuly Get product Detail",
		Data: *detailProduct,
	}
	
	ctx.JSON(http.StatusOK, response)
}

func (pc *productController) UpdateProduct(ctx *gin.Context) {
	convertedId, err := strconv.Atoi(ctx.Param("id"))
	
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert id. Make sure id is number", "details": err.Error()})
	  return
	}
  
	product := entity.Product{}
  
	detailProduct, err := pc.productRepository.GetDetailProduct(convertedId, &product)
	if err != nil {
	  ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get product data", "details": err.Error()})
	  return
	}

	updateProduct := entity.Product{}
  
	err = ctx.ShouldBind(&updateProduct)
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Input", "details": err.Error()})
	  return
	}
  
	// Update existing product data
	if strings.TrimSpace(updateProduct.Product_name) != "" {
	  detailProduct.Product_name = updateProduct.Product_name
	}
	if updateProduct.Price != 0 {
	  detailProduct.Price = updateProduct.Price
	}
	if strings.TrimSpace(updateProduct.Unit) != "" {
	  detailProduct.Unit = updateProduct.Unit
	}
  
	updatedProduct, err := pc.productRepository.UpdateProduct(convertedId,detailProduct) // Assuming updateProduct function exists
	if err != nil {
	  ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product data", "details": err.Error()})
	  return
	}
  
	// Create The response struct 
	response := ProductResponse{
	  Message: "Successfully Updated Product Data",
	  Data: *updatedProduct,
	}
  
	ctx.JSON(http.StatusOK, response) // Updated status code to 200 (OK)
  }

  func (pc *productController) DeleteProduct(ctx *gin.Context) {
	convertedId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
	  ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed convert id. Make sure id is number", "details": err.Error()})
	  return
	}

	product := entity.Product{}

	isProductExist,err := pc.productRepository.IsProductExist(convertedId,&product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error While Checking Product", "details" : err.Error()})
		return
	}
	if !isProductExist {
		ctx.JSON(http.StatusNotFound, gin.H{"message" : "product not found"})
		return
	}

	transaction_detail := entity.Transaction_detail{}

	isProductInTransaction,err := pc.productRepository.ProductInTransactionDetail(convertedId,&transaction_detail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,gin.H{"message" : "Error While Checking Product in Transaction", "details" : err.Error()})
		return
	} 
	if isProductInTransaction {
		ctx.JSON(http.StatusConflict,gin.H{"message" : "Product is being used in transaction. Please delete the transaction first"})
		return
	}

	_,err = pc.productRepository.DeleteProduct(convertedId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message" : "Error Deleting Product", "details" :err.Error()})
		return
	}

	response := struct {
		Message string `json:"message"`
		Data string `json:"data"`
	}{
		Message: "Successfully deleted data",
		Data: "OK",
	}

	ctx.JSON(http.StatusOK,response)
}