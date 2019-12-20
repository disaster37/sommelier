package repository

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/sommelier/sommelier/address"
	"github.com/sommelier/sommelier/errors"
	"github.com/sommelier/sommelier/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/bson"
)

type mongoWineEstateRepository struct {
	Conn *mongo.Collection
}

// NewMongoWineEstateRepository will create mongo handler that represent the wineEstate.Repository interface
func NewMongoWineEstateRepository(db *mongo.Database, collection string) (address.Repository, error) {

	// Check ans logs parameters
	if db == nil {
		return nil, errors.NewAPIError(errors.CodeBadParameter, "DB can't be null")
	}
	if collection == "" {
		return nil, errors.NewAPIError(errors.CodeBadParameter, "Collection can be empty")
	}
	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Database name: %s", db.Name())
		log.Debugf("Collection name : %s", collection)
	}

	return &mongoWineEstateRepository{
		Conn: db.Collection(collection),
	}, nil
}

// Fetch permit to get all documents on wine_estate collection
func (h *mongoWineEstateRepository) Fetch(ctx context.Context) ([]*models.Address, error) {

	filter := bson.M{}
	nbDocuments, err := h.Conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	wineEstateList := make([]models.Address, 0, nbDocuments)
	cursor, err := h.Conn.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		wineEstate := new(models.WineEstate)
		if err = cursor.Decode(wineEstate); err != nil {
			return nil, err
		}
		wineEstateList = append(wineEstateList, wineEstate)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Wine estates: %+v")
	}
	return wineEstateList, nil

}

// GetByID permit to get one document by ID on wine_estate collection
func (h *mongoWineEstateRepository) GetByID(ctx context.Context, id string) (*models.WineEstate, error) {

	// Get and log parameters
	if id == "" {
		return nil, errors.NewAPIError(errors.CodeBadParameter, "Wine estate ID can be empty")
	}
	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Wine estate ID: %s", id)
	}

	filter := bson.D{{"_id", id}}
	wineEstate := new(models.WineEstate)
	err := h.Conn.FindOne(
		ctx,
		filter,
	).Decode(wineEstate)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.NewAPIErrorWithError(errors.CodeDBError, err)
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Wine estate: %s", wineEstate)
	}
	return wineEstate
}

// Create permit to add new document on wine_estate collection
func (h *mongoWineEstateRepository) Create(ctx context.Context, m *models.WineEstate) (string, error) {

	// Check and logs parameters
	if m == nil {
		return "", errors.NewAPIError(errors.CodeBadParameter, "Wine estate can't be null")
	}
	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Wine estate: %s", m)
	}

	res, err := m.Conn.InsertOne(ctx, m)
	if err != nil {
		return "", errors.NewAPIErrorWithError(errors.CodeDBError, err)
	}
	id := res.InsertedID.(string)

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Wine estate ID: %s", id)
	}

	return id, nil
}

// Update permit to update document on wine_estate collection
func (h *mongoWineEstateRepository) Update(ctx context.Context, m *models.WineEstate) error {

	if m == nil {
		return errors.NewAPIError(errors.CodeBadParameter, "Wine estate can't be null")
	}
	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Wine estate: %s", m)
	}

	filter := bson.D{{"_id", a.ID}}
	_, err := m.Conn.ReplaceOne(
		ctx,
		filter,
		m,
	)

	if err != nil {
		return errors.NewAPIErrorWithError(errors.CodeDBError, err)
	}

	return nil
}

func (h *mongoWineEstateRepository) Delete(ctx context.Context, id string) error {

	if id == "" {
		return errors.NewAPIError(errors.CodeBadParameter, "Wine estate ID can't be empty")
	}
	if log.GetLevel() == log.DebugLevel {
		log.Debugf("Address ID: %s", id)
	}

	deleteDocument := new(bson.M)
	err := m.Conn.FindOneAndDelete(
		ctx,
		bson.D{{"_id", id}},
	).Decode(&deleteDocument)

	if err != nil {
		return errors.NewAPIErrorWithError(errors.CodeDBError, err)
	}

	return nil

}
