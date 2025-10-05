package repo

import (
	"fmt"
	"os"

	"github.com/Cheemx/learning-go-grpc-graphql/internal/entities"
	"github.com/Cheemx/learning-go-grpc-graphql/protobuf/golang_protobuf_brand"
	"google.golang.org/protobuf/proto"
)

type BrandRepo struct {
	brands []entities.Brand
}

const STORAGE_FILE = "./brands-storage.pb"

func (b *BrandRepo) saveToFileStorage() error {

	brandsMessage := &golang_protobuf_brand.ProtoBrandRepo{
		Brands: []*golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{},
	}

	for _, b := range b.brands {
		brandsMessage.Brands = append(brandsMessage.Brands, &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{
			ID:   uint64(b.ID),
			Name: b.Name,
			Year: uint32(b.Year),
		})
	}

	data, err := proto.Marshal(brandsMessage)
	if err != nil {
		return fmt.Errorf("cannot marshal to binary: %v", err)
	}

	err = os.WriteFile(STORAGE_FILE, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %v", err)
	}

	return nil
}

func (b *BrandRepo) loadFromFileStorage() error {
	_, err := os.Stat(STORAGE_FILE)
	if err != nil {
		fmt.Println("storage file is not found, starting with empty storage")
		return nil
	}

	data, err := os.ReadFile(STORAGE_FILE)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %v", err)
	}

	var brandsMessage golang_protobuf_brand.ProtoBrandRepo
	err = proto.Unmarshal(data, &brandsMessage)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary data to protobuf: %v", err)
	}

	for _, brand := range brandsMessage.Brands {
		b.brands = append(b.brands, entities.Brand{
			ID:   uint(brand.ID),
			Name: brand.Name,
			Year: uint(brand.Year),
		})
	}

	return nil
}

func NewBrandRepo() *BrandRepo {
	var br = BrandRepo{make([]entities.Brand, 0)}
	br.loadFromFileStorage()
	return &br
}

func (b *BrandRepo) Create(partial entities.Brand) entities.Brand {
	newItem := entities.Brand{
		ID:   uint(len(b.brands)) + 1,
		Name: partial.Name,
		Year: partial.Year,
	}
	b.brands = append(b.brands, newItem)
	b.saveToFileStorage()
	return newItem
}

func (b *BrandRepo) GetList() []entities.Brand {
	return b.brands
}

func (p *BrandRepo) GetOne(id uint) (entities.Brand, error) {
	for _, it := range p.brands {
		if it.ID == id {
			return it, nil
		}
	}
	return entities.Brand{}, fmt.Errorf("key '%d' not found", id)
}

func (p *BrandRepo) Update(id uint, amended entities.Brand) (entities.Brand, error) {
	for i, it := range p.brands {
		if it.ID == id {
			amended.ID = id
			p.brands = append(p.brands[:i], p.brands[i+1:]...)
			p.brands = append(p.brands, amended)
			return amended, nil
		}
	}
	p.saveToFileStorage()
	return entities.Brand{}, fmt.Errorf("key '%d' not found", amended.ID)
}

func (p *BrandRepo) DeleteOne(id uint) (bool, error) {
	for i, it := range p.brands {
		if it.ID == id {
			p.brands = append(p.brands[:i], p.brands[i+1:]...)
			return true, nil
		}
	}
	p.saveToFileStorage()
	return false, fmt.Errorf("key '%d' not found", id)
}

func ToProtoBrand(b entities.Brand) *golang_protobuf_brand.ProtoBrandRepo_ProtoBrand {
	return &golang_protobuf_brand.ProtoBrandRepo_ProtoBrand{
		ID:   uint64(b.ID),
		Name: b.Name,
		Year: uint32(b.Year),
	}
}

func ToBrand(b *golang_protobuf_brand.ProtoBrandRepo_ProtoBrand) entities.Brand {
	return entities.Brand{
		ID:   uint(b.ID),
		Name: b.Name,
		Year: uint(b.Year),
	}
}
