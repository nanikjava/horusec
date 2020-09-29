// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analysis

import (
	"encoding/json"
	"io"
	"time"

	apiEntities "github.com/ZupIT/horusec/development-kit/pkg/entities/api"

	horusecEntities "github.com/ZupIT/horusec/development-kit/pkg/entities/horusec"
	EnumErrors "github.com/ZupIT/horusec/development-kit/pkg/enums/errors"
	"github.com/ZupIT/horusec/development-kit/pkg/enums/horusec"
	"github.com/ZupIT/horusec/development-kit/pkg/enums/languages"
	"github.com/ZupIT/horusec/development-kit/pkg/enums/severity"
	"github.com/ZupIT/horusec/development-kit/pkg/enums/tools"
	brokerPacket "github.com/ZupIT/horusec/development-kit/pkg/services/broker/packet"
	jsonUtils "github.com/ZupIT/horusec/development-kit/pkg/utils/json"
	"github.com/ZupIT/horusec/development-kit/pkg/utils/logger"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type Interface interface {
	NewAnalysisRunning() *horusecEntities.Analysis
	DecodeAnalysisDataFromIoRead(body io.ReadCloser) (analysisData *apiEntities.AnalysisData, err error)
	ParseInterfaceToListAnalysis(data interface{}) (analysis []*horusecEntities.Analysis, err error)
	ParseInterfaceToAnalysis(data interface{}) (analysis *horusecEntities.Analysis, err error)
	ParsePacketToAnalysis(packet brokerPacket.IPacket) (analysis *horusecEntities.Analysis, err error)
	SetFindOneFilter(analysisID string) map[string]interface{}
}

type UseCases struct {
}

func NewAnalysisUseCases() Interface {
	return &UseCases{}
}

func (au *UseCases) NewAnalysisRunning() *horusecEntities.Analysis {
	return &horusecEntities.Analysis{
		ID:              uuid.New(),
		Status:          horusec.Running,
		Errors:          "",
		CreatedAt:       time.Now(),
		Vulnerabilities: []horusecEntities.Vulnerability{},
	}
}

func (au *UseCases) ParsePacketToAnalysis(
	packet brokerPacket.IPacket) (analysis *horusecEntities.Analysis, err error) {
	if err = json.Unmarshal(packet.GetBody(), &analysis); err != nil {
		logger.LogError(EnumErrors.ErrParsePacketToAnalysis, err)
		return analysis, err
	}

	return analysis, nil
}

func (au *UseCases) ParseInterfaceToListAnalysis(data interface{}) (analysis []*horusecEntities.Analysis, err error) {
	return analysis, jsonUtils.ConvertInterfaceToOutput(data, &analysis)
}

func (au *UseCases) ParseInterfaceToAnalysis(data interface{}) (*horusecEntities.Analysis, error) {
	analysis := &horusecEntities.Analysis{}
	return analysis, jsonUtils.ConvertInterfaceToOutput(data, analysis)
}

func (au *UseCases) SetFindOneFilter(analysisID string) map[string]interface{} {
	return map[string]interface{}{"id": analysisID}
}

func (au *UseCases) DecodeAnalysisDataFromIoRead(body io.ReadCloser) (
	analysisData *apiEntities.AnalysisData, err error) {
	if body == nil {
		return nil, EnumErrors.ErrorBodyIsRequired
	}
	err = json.NewDecoder(body).Decode(&analysisData)
	_ = body.Close()
	if err != nil {
		return nil, err
	}
	return analysisData, au.validateAnalysis(analysisData.Analysis)
}

func (au *UseCases) validateAnalysis(analysis *horusecEntities.Analysis) error {
	return validation.ValidateStruct(analysis,
		validation.Field(&analysis.ID, validation.Required, is.UUID),
		validation.Field(&analysis.Status,
			validation.Required, validation.In(horusec.Running, horusec.Success, horusec.Error)),
		validation.Field(&analysis.CreatedAt, validation.Required, validation.NilOrNotEmpty),
		validation.Field(&analysis.FinishedAt, validation.Required, validation.NilOrNotEmpty),
		validation.Field(&analysis.Vulnerabilities, validation.By(au.validateVulnerabilities(analysis.Vulnerabilities))),
	)
}

func (au *UseCases) validateVulnerabilities(vulnerabilities []horusecEntities.Vulnerability) validation.RuleFunc {
	return func(value interface{}) error {
		if len(vulnerabilities) == 0 {
			return nil
		}
		for key := range vulnerabilities {
			if err := validation.ValidateStruct(&vulnerabilities[key],
				validation.Field(&vulnerabilities[key].SecurityTool, validation.Required, validation.In(au.sliceTools()...)),
				validation.Field(&vulnerabilities[key].Language, validation.Required, validation.In(au.sliceLanguages()...)),
				validation.Field(&vulnerabilities[key].Severity, validation.Required, validation.In(au.sliceSeverities()...)),
			); err != nil {
				return err
			}
		}
		return nil
	}
}

func (au *UseCases) sliceTools() []interface{} {
	return []interface{}{
		tools.GoSec,
		tools.SecurityCodeScan,
		tools.GitLeaks,
		tools.Brakeman,
		tools.NpmAudit,
		tools.Safety,
		tools.Bandit,
		tools.SpotBugs,
		tools.YarnAudit,
		tools.TfSec,
		tools.HorusecJava,
		tools.HorusecKotlin,
		tools.HorusecLeaks,
	}
}
func (au *UseCases) sliceLanguages() []interface{} {
	return []interface{}{
		languages.Go,
		languages.DotNet,
		languages.Ruby,
		languages.Python,
		languages.Java,
		languages.Kotlin,
		languages.Javascript,
		languages.Leaks,
		languages.HCL,
	}
}
func (au *UseCases) sliceSeverities() []interface{} {
	return []interface{}{
		severity.Info,
		severity.NoSec,
		severity.Low,
		severity.Medium,
		severity.High,
		severity.Audit,
	}
}