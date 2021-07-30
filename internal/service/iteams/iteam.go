package routes

import (
	"github.com/DrIhor/test_task/internal/models/iteams"
	iteamModel "github.com/DrIhor/test_task/internal/models/iteams"
)

type IteamServices struct {
	storage iteamModel.IteamStorageServices
}

func New(stor iteamModel.IteamStorageServices) *IteamServices {
	return &IteamServices{
		storage: stor,
	}
}

func (iteamSrv *IteamServices) AddNewIteam(iteam iteams.Iteam) error {
	return iteamSrv.storage.AddNewIteam(iteam)
}

func (iteamSrv *IteamServices) GetIteam(iteamName string) ([]byte, error) {
	return iteamSrv.storage.GetIteam(iteamName)
}

func (iteamSrv *IteamServices) DeleteIteam(iteamName string) error {
	return iteamSrv.storage.DeleteIteam(iteamName)
}

func (iteamSrv *IteamServices) UpdateIteam(iteamName string) ([]byte, error) {
	return iteamSrv.storage.UpdateIteam(iteamName)
}
