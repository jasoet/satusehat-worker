package model

type DiagnosisList []Diagnosis

func (dl *DiagnosisList) Invalid() bool {
	if dl == nil {
		return true
	}
	for _, diagnosis := range *dl {
		if diagnosis.Invalid() {
			return true
		}
	}
	return false
}

type MedicationRequestList []MedicationRequest

func (mrl *MedicationRequestList) Invalid() bool {
	if mrl == nil {
		return true
	}
	for _, request := range *mrl {
		if request.Invalid() {
			return true
		}
	}
	return false
}

type MedicationDispenseList []MedicationDispense

func (mdl *MedicationDispenseList) Invalid() bool {
	if mdl == nil {
		return true
	}
	for _, dispense := range *mdl {
		if dispense.Invalid() {
			return true
		}
	}
	return false
}

type ProcedureList []Procedure

func (pl *ProcedureList) Invalid() bool {
	if pl == nil {
		return true
	}
	for _, procedure := range *pl {
		if procedure.Invalid() {
			return true
		}
	}
	return false
}

type ObservationLabList []ObservationLab

func (oll *ObservationLabList) Invalid() bool {
	if oll == nil {
		return true
	}
	for _, lab := range *oll {
		if lab.Invalid() {
			return true
		}
	}
	return false
}

type ObservationRadiologyList []ObservationRadiology

func (orl *ObservationRadiologyList) Invalid() bool {
	if orl == nil {
		return true
	}
	for _, radiology := range *orl {
		if radiology.Invalid() {
			return true
		}
	}
	return false
}
