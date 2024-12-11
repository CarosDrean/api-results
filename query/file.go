package query

import "github.com/CarosDrean/api-results.git/models"

var ExcelFile = models.QueryDB{
	"getDataFile": {Q: `select per.v_PersonId, ser.v_ServiceId, CONCAT(per.v_FirstLastName, ' ' ,per.v_SecondLastName, ' ', per.v_FirstName) V_NombrePersona, per.v_FirstLastName, per.v_SecondLastName, per.v_FirstName, per.v_DocNumber, per.i_SexTypeId, 
		per.v_BirthPlace ,per.v_AdressLocation, per.d_Birthdate, syspara.v_Value1, org.v_Name as v_OrgName, ser.d_GlobalExpirationDate, loc.v_Name as v_Location,
		proto.v_Name, ser.d_ServiceDate, per.v_CurrentOccupation, (select syspad.v_Value1 from systemparameter syspad where syspad.i_ParameterId = ser.i_AptitudeStatusId  and syspad.i_GroupId = '124') V_Aptitude
		from service ser join person per on ser.v_PersonId = per.v_PersonId
		join calendar ca on ser.v_ServiceId = ca.v_ServiceId 
		join protocol proto on ser.v_ProtocolId = proto.v_ProtocolId
		join systemparameter syspara on proto.i_EsoTypeId = syspara.i_ParameterId
		join organization org on proto.v_CustomerOrganizationId = org.v_OrganizationId
		join location loc on proto.v_CustomerLocationId = loc.v_LocationId
		where ser.d_ServiceDate >= CONCAT('%s', ' 00:00:00.0000000') and  ser.d_ServiceDate <= CONCAT('%s', ' 23:59:59.9999999') and syspara.i_GroupId = '118' and proto.i_EsoTypeId != 4 and proto.i_EsoTypeId != 6 
		and ca.i_CalendarStatusId != 4 and proto.v_EmployerOrganizationId = '%s' order by ser.d_ServiceDate desc;`},

	"getInterconsultas": {Q: `select mres.v_Name, res.v_ServiceId, res.v_DiagnosticRepositoryId from restriction res 
		join masterrecommendationrestricction mres on res.v_MasterRestrictionId = mres.v_MasterRecommendationRestricctionId 
		where res.i_IsDeleted = '0' and mres.i_TypifyingId = '2' and res.v_ServiceId = '%s';`},

	"getRestriccioens": {Q: `select ISNULL(mres.v_Name,'---') v_Name from restriction res 
		join masterrecommendationrestricction mres on res.v_MasterRestrictionId = mres.v_MasterRecommendationRestricctionId 
		where res.i_IsDeleted = '0' and mres.i_TypifyingId = '2' and res.v_ServiceId = '%s';`},

	"getRecomendaciones": {Q: `select mresc.v_Name from recommendation rec
		join masterrecommendationrestricction mresc on rec.v_MasterRecommendationId = mresc.v_MasterRecommendationRestricctionId
		where rec.v_DiagnosticRepositoryId = '%s' and rec.v_ServiceId = '%s'
		and rec.i_IsDeleted = '0';`},

	"getAptitudAltura": {Q: `select sys.v_Value1 from service ser 
		join servicecomponent serc on ser.v_ServiceId = serc.v_ServiceId
		join servicecomponentfields serf on serc.v_ServiceComponentId = serf.v_ServiceComponentId
		join servicecomponentfieldvalues serv on serf.v_ServiceComponentFieldsId = serv.v_ServiceComponentFieldsId
		join component com on serc.v_ComponentId = com.v_ComponentId
		join componentfields comfs on serf.v_ComponentFieldId = comfs.v_ComponentFieldId
		join componentfield comf on serf.v_ComponentFieldId = comf.v_ComponentFieldId
		join component comp on comfs.v_ComponentId = comp.v_ComponentId
		join systemparameter sys on serv.v_Value1 = sys.i_ParameterId
		where ser.v_ServiceId = '%s' and comp.v_ComponentId = 'N009-ME000000015' 
		and comf.v_ComponentFieldId = 'N009-MF000000039' and serc.i_IsDeleted = 0 and serf.i_IsDeleted = 0
		and sys.i_GroupId = 163;`},

	"getAptitudEspacios": {Q: `select serv.v_Value1 from service ser 
		join servicecomponent serc on ser.v_ServiceId = serc.v_ServiceId
		join servicecomponentfields serf on serc.v_ServiceComponentId = serf.v_ServiceComponentId
		join servicecomponentfieldvalues serv on serf.v_ServiceComponentFieldsId = serv.v_ServiceComponentFieldsId
		join component com on serc.v_ComponentId = com.v_ComponentId
		join componentfields comfs on serf.v_ComponentFieldId = comfs.v_ComponentFieldId
		join componentfield comf on serf.v_ComponentFieldId = comf.v_ComponentFieldId
		join component comp on comfs.v_ComponentId = comp.v_ComponentId
		where ser.v_ServiceId = '%s' and comp.v_ComponentId = 'N009-ME000000436' 
		and comf.v_ComponentFieldId = 'N009-MF000003359' and serc.i_IsDeleted = 0 and serf.i_IsDeleted = 0;`},

	"getCustomeValueV1": {Q: `select serv.v_Value1 from service ser 
		join servicecomponent serc on ser.v_ServiceId = serc.v_ServiceId
		join servicecomponentfields serf on serc.v_ServiceComponentId = serf.v_ServiceComponentId
		join servicecomponentfieldvalues serv on serf.v_ServiceComponentFieldsId = serv.v_ServiceComponentFieldsId
		join component com on serc.v_ComponentId = com.v_ComponentId
		join componentfields comfs on serf.v_ComponentFieldId = comfs.v_ComponentFieldId
		join componentfield comf on serf.v_ComponentFieldId = comf.v_ComponentFieldId
		join component comp on comfs.v_ComponentId = comp.v_ComponentId
		where ser.v_ServiceId = '%s' and comp.v_ComponentId = '%s' 
		and comf.v_ComponentFieldId = '%s' and serc.i_IsDeleted = 0 and serf.i_IsDeleted = 0;`},

	"getCustomeValueV2": {Q: `select serv.v_Value1 from service ser 
		join servicecomponent serc on ser.v_ServiceId = serc.v_ServiceId
		join servicecomponentfields serf on serc.v_ServiceComponentId = serf.v_ServiceComponentId
		join servicecomponentfieldvalues serv on serf.v_ServiceComponentFieldsId = serv.v_ServiceComponentFieldsId
		where ser.v_ServiceId = '%s' and serc.v_ComponentId = '%s' 
		and serf.v_ComponentFieldId = '%s' and serc.i_IsDeleted = 0 and serf.i_IsDeleted = 0;`},

	"getValueFromParametersV1": {Q: `select v_Value1 from systemparameter
		where i_GroupId = '%s' and i_ParameterId = '%s';`},

	"getDxSingle": {Q: `select d.v_Name from diagnosticrepository dr inner join diseases d on dr.v_DiseasesId = d.v_DiseasesId
  		where dr.v_ServiceId = '%s' and dr.v_ComponentId = '%s';`},

	"getCheckDx": {Q: `select d.v_Name from diagnosticrepository dr inner join diseases d on dr.v_DiseasesId = d.v_DiseasesId
  		where dr.v_ServiceId = '%s' and dr.v_DiseasesId = '%s';`},

	"getNoxiusHabitats": {Q: `select nh.i_TypeHabitsId, sp.v_Value1, nh.v_Frequency, nh.v_Comment from calendar c 
		inner join noxioushabits nh on c.v_PersonId = nh.v_PersonId
		inner join systemparameter sp on nh.i_TypeHabitsId = sp.i_ParameterId
		where c.v_ServiceId = '%s' and sp.i_GroupId = '148' and nh.i_TypeHabitsId = '%s';`},

	"getAntecedentesPersonales": {Q: `select v_DiagnosticDetail from personmedicalhistory where v_PersonId = '%s' and v_DiseasesId != 'N009-DD000001698';`},

	"getCheckAntePerso": {Q: `select v_DiagnosticDetail from personmedicalhistory where v_PersonId = '%s' and v_DiseasesId = '%s';`},
}
