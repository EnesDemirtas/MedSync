package medicinedb

import (
	"database/sql"
	"time"

	"github.com/EnesDemirtas/medisync/business/core/crud/medicine"
	"github.com/EnesDemirtas/medisync/business/data/sqldb/dbarray"
	"github.com/google/uuid"
)

type dbMedicine struct {
	ID 			 uuid.UUID		`db:"medicine_id"`
	Name		 string			`db:"name"`
	Description  sql.NullString	`db:"description"`
	Manufacturer sql.NullString	`db:"manufacturer"`
	Type		 sql.NullString `db:"type"`
	Tags		 dbarray.String	`db:"tags"`
	ExpiryDate	 time.Time		`db:"expiry_date"`
	DateCreated  time.Time		`db:"date_created"`
	DateUpdated  time.Time		`db:"date_updated"`
}

func toDBMedicine(med medicine.Medicine) dbMedicine {
	tags := make([]string, len(med.Tags))
	for i, tag := range med.Tags {
		tags[i] = tag.String()
	}

	return dbMedicine{
		ID:			  med.ID,
		Name:		  med.Name,
		Description:  sql.NullString{
			String: med.Description,
			Valid:	med.Description != "",
		},
		Manufacturer: sql.NullString{
			String: med.Manufacturer,
			Valid:  med.Manufacturer != "",
		},
		Type: 		  sql.NullString{
			String:	med.Type,
			Valid:  med.Type != "",
		},
		Tags: 		  tags,
		ExpiryDate:   med.ExpiryDate,
		DateCreated:  med.DateCreated,
		DateUpdated:  med.DateUpdated,
	}
}

func toCoreMedicine(dbMedicine dbMedicine) medicine.Medicine {
	tags := make([]uuid.UUID, len(dbMedicine.Tags))
	for i, tag := range dbMedicine.Tags {
		tags[i] = uuid.MustParse(tag)
	}
	med := medicine.Medicine{
		ID:			  dbMedicine.ID,
		Name:		  dbMedicine.Name,
		Description:  dbMedicine.Description.String,
		Manufacturer: dbMedicine.Manufacturer.String,
		Type:		  dbMedicine.Type.String,
		Tags:		  tags,
		ExpiryDate:   dbMedicine.ExpiryDate,
		DateCreated:  dbMedicine.DateCreated,
		DateUpdated:  dbMedicine.DateUpdated,
	}

	return med
}

func toCoreMedicineSlice(dbMedicines []dbMedicine) []medicine.Medicine {
	meds := make([]medicine.Medicine, len(dbMedicines))

	for i, dbMedicine := range dbMedicines {
		meds[i] = toCoreMedicine(dbMedicine)
	}

	return meds
}