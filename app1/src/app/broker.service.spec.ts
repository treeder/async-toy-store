import { TestBed } from '@angular/core/testing';

import { BrokerService } from './broker.service';

describe('BrokerService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: BrokerService = TestBed.get(BrokerService);
    expect(service).toBeTruthy();
  });
});
