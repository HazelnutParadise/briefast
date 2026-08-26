## ADDED Requirements

### Requirement: Advertising slot

Report pages SHALL carry one advertising slot placed between the stock calls section and the industry news section. The slot SHALL be visually delimited from editorial content by a rule and a label reading 廣告, because the advertisement itself is rendered by a third party whose appearance the application does not control.

The slot SHALL consist of a mount point rendered as part of the page markup and a loader script delivered through a node that executes scripts in the main document. The loader script SHALL NOT re-execute when the page re-renders in response to a background update, so that a refresh of the report does not count an additional impression.

Pages without report content SHALL NOT carry the slot.

#### Scenario: Slot appears between calls and industries

- **WHEN** a client requests a page showing a report
- **THEN** the rendered markup contains the advertising mount point after the stock calls section and before the industry news section

#### Scenario: Slot is labelled and delimited

- **WHEN** the advertising slot is rendered
- **THEN** it carries a label reading 廣告 and a rule separating it from the surrounding editorial content

#### Scenario: Loader runs in the main document

- **WHEN** a report page is rendered in a browser
- **THEN** the loader script executes and resolves the mount point that the page markup provided

#### Scenario: Background update does not re-run the loader

- **WHEN** a new report is published and the open page re-renders through the live update channel
- **THEN** the advertising node is reused and its loader script does not run again

#### Scenario: Pages without a report carry no slot

- **WHEN** a client requests the history listing, or a dated view for which no report exists
- **THEN** the rendered markup contains no advertising mount point

##### Example: slot presence by page

| Page | Advertising slot |
|---|---|
| home page showing the latest report | present |
| dated view of a stored report | present |
| history listing | absent |
| dated view with no stored report | absent |
| waiting-for-report home page | absent |
