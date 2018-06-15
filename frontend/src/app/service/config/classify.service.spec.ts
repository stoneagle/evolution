import { TestBed, inject } from '@angular/core/testing';

import { ClassifyService } from './classify.service';

describe('ClassifyService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ClassifyService]
    });
  });

  it('should be created', inject([ClassifyService], (service: ClassifyService) => {
    expect(service).toBeTruthy();
  }));
});
