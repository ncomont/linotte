package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	"git.ve.home/nicolasc/linotte/job_runner/configuration"
	"git.ve.home/nicolasc/linotte/libs/csv"
	"git.ve.home/nicolasc/linotte/libs/helpers"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
)

type entry struct {
	CD_NOM   uint32
	NAME_LAT string
	NAME_FR  string
	CAT_LRR  string
	NAME_POP string
	COM      string
}

// RedListRunner represents the RedList specific runner
type RedListRunner struct {
	taxref   proto_taxref.TaxrefClient
	job      *proto_job.JobReply
	channels *Channels
}

// NewRedListRunner creates a new redlist runner
func NewRedListRunner() Runner {
	connection, err := grpc.Dial(config.Get().TaxrefServiceEndpoint, grpc.WithInsecure())
	// FIXME: when endpoint is invalid, no error seems to be thrown
	helpers.HandleError(err)

	return RedListRunner{
		taxref: proto_taxref.NewTaxrefClient(connection),
	}
}

// Configure helps to configure the runner. It accepts the output channels
func (runner RedListRunner) Configure(channels *Channels) Runner {
	runner.channels = channels
	return runner
}

// Run starts the RedList runner with the given job
func (runner RedListRunner) Run(job *proto_job.JobReply) {
	var (
		lines   chan csv.Line = make(chan csv.Line)
		entries []entry
		results []*proto_job.ResultReply
		err     error
		element entry
	)

	runner.job = job

	path := fmt.Sprintf(
		"%s/%s/%s", config.Get().StoragePath, runner.job.Type, runner.job.File,
	)
	fmt.Printf("Running %s\n", path)

	go csv.ReadFile(path, lines, ';')
	for msg := range lines {
		if msg.Error != nil {
			runner.channels.Errors <- msg.Error
			return
		}

		if element, err = parse(msg.Elements); err != nil {
			runner.channels.Errors <- err
			return
		}

		entries = append(entries, element)
	}

	// TODO: don't wait the end of the file read to begin the job ?

	if len(entries) == 0 {
		runner.channels.Errors <- errors.New("No entries for list")
	} else {
		if results, err = runner.process(entries); err != nil {
			runner.channels.Errors <- err
		} else {
			runner.channels.Results <- results
		}
	}
}

func (runner RedListRunner) process(entries []entry) ([]*proto_job.ResultReply, error) {
	var (
		results []*proto_job.ResultReply
		id      uint32
		err     error
	)

	for _, e := range entries {
		if id, err = runner.findTaxonID(e); err != nil {
			return nil, err
		}

		serialized, err := json.Marshal(e)
		if err != nil {
			return nil, err
		}

		result := proto_job.ResultReply{
			SearchData: string(serialized[:]),
			Value:      e.CAT_LRR,
		}

		if id > 0 {
			result.State = "FOUND"
			result.TaxonId = id
		} else {
			result.State = "NOT_FOUND"
		}

		results = append(results, &result)
		runner.channels.Progression <- uint32(float32(len(results)) / float32(len(entries)) * 100)
	}

	return results, nil
}

func (runner RedListRunner) findTaxonID(e entry) (uint32, error) {
	var (
		taxon = &proto_taxref.TaxonReply{}
		err   error
	)

	if e.CD_NOM > 0 {
		if taxon, err = runner.taxref.ReferenceByID(
			context.Background(),
			&proto_taxref.TaxonRequest{
				Id: e.CD_NOM,
			}); err != nil {
			return 0, err
		}
	}

	if taxon.Id == 0 && len(e.NAME_LAT) > 0 {
		if taxon, err = runner.taxref.ReferenceByVerb(
			context.Background(),
			&proto_taxref.TaxonRequest{
				Name: e.NAME_LAT,
			}); err != nil {
			return 0, err
		}

		if taxon.Id == 0 {
			if taxon, err = runner.taxref.ReferenceByVerb(
				context.Background(),
				&proto_taxref.TaxonRequest{
					FullName:          e.NAME_LAT,
					IgnorePunctuation: true,
				}); err != nil {
				return 0, err
			}
		}
	}

	if taxon.Id == 0 && len(e.NAME_FR) > 0 {
		if taxon, err = runner.taxref.ReferenceByVerb(
			context.Background(),
			&proto_taxref.TaxonRequest{
				VernacularName: e.NAME_FR,
			}); err != nil {
			return 0, err
		}
	}

	if taxon.Id == 0 {
		return runner.substitute(e.NAME_LAT)
	}

	return taxon.Id, nil
}

func (runner RedListRunner) substitute(name string) (uint32, error) {
	var (
		taxon  = &proto_taxref.TaxonReply{}
		err    error
		result []string
	)

	name = strings.Replace(name, "L.", " ", -1)
	name = strings.Replace(name, " ET ", " ", -1)
	name = strings.Replace(name, " et ", " ", -1)
	name = strings.Replace(name, " & ", " ", -1)
	name = strings.Replace(name, " ex ", " ", -1)
	name = regexp.MustCompile(`\([^)]*\)`).ReplaceAllString(name, " ")
	name = regexp.MustCompile(`[0-9]+`).ReplaceAllString(name, " ")

	parts := strings.Split(name, " ")
	result = append(result, parts[0])
	for i, w := range parts {
		startsWithUpper, _ := regexp.MatchString(`^[A-Z].*`, w)
		if i > 0 && w != "" && w != " " && !startsWithUpper {
			result = append(result, w)
		}
	}

	term := strings.Join(result, " ")
	if taxon, err = runner.taxref.ReferenceByVerb(
		context.Background(),
		&proto_taxref.TaxonRequest{
			Name: term,
		}); err != nil {
		return 0, err
	}

	if taxon.Id == 0 {
		if taxon, err = runner.taxref.ReferenceByVerb(
			context.Background(),
			&proto_taxref.TaxonRequest{
				FullName:          term,
				IgnorePunctuation: true,
			}); err != nil {
			return 0, err
		}
	}

	return taxon.Id, nil
}

func parse(elements []string) (entry, error) {
	var (
		err error
		id  uint64
	)

	if len(elements) < 6 {
		return entry{}, errors.New("Invalid file: wrong columns count")
	}

	if elements[0] != "" {
		id, err = strconv.ParseUint(elements[0], 10, 64)
		helpers.HandleError(err)
	}

	return entry{
		CD_NOM:   uint32(id),
		NAME_LAT: elements[1],
		NAME_FR:  elements[2],
		CAT_LRR:  elements[3],
		NAME_POP: elements[4],
		COM:      elements[5],
	}, nil
}
