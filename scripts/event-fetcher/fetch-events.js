const fs = require('fs');

const queryComposition = (compositionId) => `query {
  compositions(condition: {id: "${compositionId}" }) {
    nodes {
      compositionVersions(orderBy: SEMVER_VERSION_DESC, first: 1) {
        nodes {
          version
          componentVersions{
            nodes {
              apiVersions (filter: {
                and: [                            
                  { governanceStatus: { equalTo: PASSED } }
                  { 
                    api: {
                      or: [
                        { apiStandard: { equalTo: SYSTEM_EVENTS_V2 } }
                        { apiStandard: { equalTo: UI_EVENTS_V1 } }
                      ]
                    } 
                  }
                ]
              }){
                nodes {
                  api {
                    name
                    id
                    apiStandard
                    entities(
                      filter: {
                        entityRevisionsExist: true
                      }
                    ){
                      nodes {
                        uniqueIdentifier
                        exposed
                        entityRevisions(orderBy: SEMVER_VERSION_DESC, first: 1) {
                          nodes {
                            stability
                            visibility
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`

const fetchQuery = (query) => {
  return fetch(process.env.APICULTURIST_BASE_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${process.env.APICULTURIST_TOKEN}`,
    },
    body: JSON.stringify({
      query,
    }),
  }).then(resp => resp.json());
}

const run = async () => {
  const jsonResponseQCScomposition = await fetchQuery(queryComposition('fe179367-18c9-4e5a-980f-139f21084dcf')); // Qlik Cloud Service id: fe179367-18c9-4e5a-980f-139f21084dcf
  const jsonResponseQMFEUIcomposition = await fetchQuery(queryComposition('8780ee5b-31ec-4d6c-a339-e1588d87acef')); // QMFE UI Events id: 8780ee5b-31ec-4d6c-a339-e1588d87acef
  let qcsEvents = [];
  let qmfeUIevents = [];
  [jsonResponseQCScomposition, jsonResponseQMFEUIcomposition].forEach((response) => {
    response?.data?.compositions?.nodes[0]?.compositionVersions?.nodes[0]?.componentVersions?.nodes.forEach((node) => {
      node.apiVersions?.nodes.forEach((apiNode) => {
        const api = apiNode.api;
        if (api) {
          if (api.apiStandard === 'SYSTEM_EVENTS_V2') {
            qcsEvents.push(api);
          } else if (api.apiStandard === 'UI_EVENTS_V1') {
            qmfeUIevents.push(api);
          }
        }
      });
    });
  });

  const events = [...qcsEvents, ...qmfeUIevents];

  const eventTypeList = events.map(api => {
    const entities = api.entities.nodes.map(entity => {
      return {
        uniqueIdentifier: entity.uniqueIdentifier,
        stability: entity.entityRevisions.nodes[0].stability,
        visibility: entity.entityRevisions.nodes[0].visibility,
        exposed: entity.exposed,
      };
    });
    return entities;
  }).flat();

  const filteredEvents = eventTypeList.filter((event) => {
    return event.visibility == 'PUBLIC' && event.stability == 'STABLE';
  });
  const targetedEvents = filteredEvents.flatMap((event) => (event.uniqueIdentifier.split("/")[1].trim()));
  const targetedEventsFiltered = targetedEvents.filter((event) => { return !event.endsWith('.purged')}).sort()
  const yamlOutput = {
    allowedEvents: targetedEventsFiltered
  };
  try {
    fs.writeFileSync('allowed-events.yaml', `allowedEvents:\n  - ${yamlOutput.allowedEvents.join('\n  - ')}`, 'utf8');
    console.log('allowed-events.yaml file has been saved.');
  } catch (err) {
    console.error('Error writing to file', err);
  }
}
run();
