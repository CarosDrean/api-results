package query

import "github.com/CarosDrean/api-results.git/models"

var Statistics = models.QueryDB{
	"getDisease": {Q: "select s." + service.Fields[0] + ", s." + service.Fields[3] + ", " +
		"p." + person.Fields[0] + ", pr." + protocol.Fields[0] + ", s." + service.Fields[6] + ", p." + person.Fields[1] + ", p." + person.Fields[3] +
		", p." + person.Fields[4] + ", p." + person.Fields[5] + ", p." + person.Fields[6] + ", p." + person.Fields[7] + ", p." + person.Fields[8] +
		", d.v_Name, c.v_Name, sp.v_Value1 from service s " +
		"inner join person p on s.v_PersonId = p.v_PersonId " +
		"left join protocol pr on s.v_ProtocolId = pr.v_ProtocolId " +
		"left join organization o on pr.v_CustomerOrganizationId = o.v_OrganizationId " +
		"inner join diagnosticrepository dr on s.v_ServiceId = dr.v_ServiceId " +
		"inner join component c on dr.v_ComponentId = c.v_ComponentId " +
		"inner join systemparameter sp on sp.i_GroupId = 116 and c.i_CategoryId = sp.i_ParameterId " +
		"inner join diseases d on dr.v_DiseasesId = d.v_DiseasesId " +
		"where dr.i_IsDeleted = 0 and s.i_ServiceStatusId = 3 and pr.v_ProtocolId = '%s' " +
		"and s.d_ServiceDate >= CAST('%s' as date) and s.d_ServiceDate <= CAST('%s' as date) " +
		"order by s.d_ServiceDate desc"},
	"getAllDiseaseDate": {Q: "select s." + service.Fields[0] + ", s." + service.Fields[3] + ", " +
		"p." + person.Fields[0] + ", pr." + protocol.Fields[0] + ", s." + service.Fields[6] + ", p." + person.Fields[1] + ", p." + person.Fields[3] +
		", p." + person.Fields[4] + ", p." + person.Fields[5] + ", p." + person.Fields[6] + ", p." + person.Fields[7] + ", p." + person.Fields[8] +
		", d.v_Name, c.v_Name, sp.v_Value1 from service s " +
		"inner join person p on s.v_PersonId = p.v_PersonId " +
		"left join protocol pr on s.v_ProtocolId = pr.v_ProtocolId " +
		"left join organization o on pr.v_CustomerOrganizationId = o.v_OrganizationId " +
		"inner join diagnosticrepository dr on s.v_ServiceId = dr.v_ServiceId " +
		"inner join component c on dr.v_ComponentId = c.v_ComponentId " +
		"inner join systemparameter sp on sp.i_GroupId = 116 and c.i_CategoryId = sp.i_ParameterId " +
		"inner join diseases d on dr.v_DiseasesId = d.v_DiseasesId " +
		"where dr.i_IsDeleted = 0 and s.i_ServiceStatusId = 3 " +
		"and s.d_ServiceDate >= CAST('%s' as date) and s.d_ServiceDate <= CAST('%s' as date) " +
		"order by s.d_ServiceDate desc"},
}
