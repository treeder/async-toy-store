import { TestBed } from '@angular/core/testing';

import { NatsService } from './nats.service';

describe('NatsService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: NatsService = TestBed.get(NatsService);
    expect(service).toBeTruthy();
  });
});
