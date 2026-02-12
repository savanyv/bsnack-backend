package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/savanyv/bsnack-backend/internal/cache"
	dtos "github.com/savanyv/bsnack-backend/internal/dto"
	"github.com/savanyv/bsnack-backend/internal/model"
	"github.com/savanyv/bsnack-backend/internal/repository"
)

type ProductUsecase interface {
	GetAll(ctx context.Context) ([]dtos.ProductResponse, error)
	CreateProduct(ctx context.Context, req dtos.CreateProductRequest) (*dtos.ProductResponse, error)
	GetByID(ctx context.Context, ID string) (*dtos.ProductResponse, error)
}

type productUsecase struct {
	productRepo repository.ProductRepository
	redis *cache.RedisClient
}

func NewProductUsecase(pr repository.ProductRepository, r *cache.RedisClient) ProductUsecase {
	return &productUsecase{
		productRepo: pr,
		redis: r,
	}
}

func (u *productUsecase) GetAll(ctx context.Context) ([]dtos.ProductResponse, error) {
	cacheKey := "products:all"
	cached, err := u.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var products []dtos.ProductResponse
		if err := json.Unmarshal([]byte(cached), &products); err != nil {
			return products, nil
		}
	}

	products, err := u.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.ProductResponse, 0)
	for _, p := range products {
		responses = append(responses, dtos.ProductResponse{
			ID: p.ID.String(),
			Name: p.Name,
			Type: p.Type,
			Flavor: p.Flavor,
			Size: p.Size,
			Price: p.Price,
			Stock: p.Stock,
			CreatedAt: p.CreatedAt,
		})
	}

	data, _ := json.Marshal(responses)
	_ = u.redis.Set(ctx, cacheKey, string(data), 10 * time.Minute)

	return responses, nil
}

func (u *productUsecase) CreateProduct(ctx context.Context, req dtos.CreateProductRequest) (*dtos.ProductResponse, error) {
	if req.Name == "" || req.Price <= 0 {
		return nil, errors.New("name and price are required")
	}

	product := &model.Product{
		Name: req.Name,
		Type: req.Type,
		Flavor: req.Flavor,
		Size: req.Size,
		Price: req.Price,
		Stock: req.Stock,
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	_ = u.redis.Delete(ctx, "products:all")

	response := dtos.ProductResponse{
		ID: product.ID.String(),
		Name: product.Name,
		Type: product.Type,
		Flavor: product.Flavor,
		Size: product.Size,
		Price: product.Price,
		Stock: product.Stock,
		CreatedAt: product.CreatedAt,
	}

	return &response, nil
}

func (u *productUsecase) GetByID(ctx context.Context, ID string) (*dtos.ProductResponse, error) {
	cacheKey := "product:" + ID
	cached, err := u.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var product dtos.ProductResponse
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	product, err := u.productRepo.FindByID(ctx, ID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	response := dtos.ProductResponse{
		ID: product.ID.String(),
		Name: product.Name,
		Type: product.Type,
		Flavor: product.Flavor,
		Size: product.Size,
		Price: product.Price,
		Stock: product.Stock,
		CreatedAt: product.CreatedAt,
	}

	data, _ := json.Marshal(response)
	_ = u.redis.Set(ctx, cacheKey, string(data), 10 * time.Minute)

	return &response, nil
}
